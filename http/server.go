package http

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/surasithaof/sse/shared"
)

type server struct {
	connections map[string]*Connection
	sync.RWMutex
}

func NewServer() SSEServer {
	s := server{
		connections: make(map[string]*Connection),
	}
	return &s
}

func (s *server) addClient(clientID string, rw http.ResponseWriter, req *http.Request) *Connection {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return nil
	}

	s.Lock()
	defer s.Unlock()
	client := &Connection{
		ID:         clientID,
		writer:     rw,
		flusher:    flusher,
		requestCtx: req.Context(),
	}
	s.connections[clientID] = client

	return client
}

func (s *server) removeClient(clientID string) {
	s.Lock()
	defer s.Unlock()
	fmt.Sprintln("Client gone")

	_, ok := s.connections[clientID]
	if ok {
		delete(s.connections, clientID)
	}
}

func (s *server) Listen(rw http.ResponseWriter, req *http.Request, connectionID string) {
	s.addClient(connectionID, rw, req)
	defer func() {
		s.removeClient(connectionID)
	}()

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	<-req.Context().Done()
}

func (s *server) SendMessage(connectionID string, event shared.Event) error {
	s.RLock()
	defer s.RUnlock()

	connection, ok := s.connections[connectionID]
	if !ok {
		// TODO: need to handle error
		return errors.New("not_found_connection")
	}
	if ok {
		err := connection.send(event.Event, event.Message)
		if err != nil {
			s.removeClient(connection.ID)
			return err
		}
	}
	return nil
}

func (s *server) Broadcast(event shared.Event) {
	s.RLock()
	defer s.RUnlock()

	for cID, connection := range s.connections {
		err := connection.send(event.Event, event.Message)
		if err != nil {
			s.removeClient(cID)
			continue
		}
	}
}
