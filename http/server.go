package http

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	connections map[string]*Connection
	sync.RWMutex
}

func NewServer() *Server {
	s := Server{}
	return &s
}

func (s *Server) addClient(clientID string, rw http.ResponseWriter, req *http.Request) *Connection {
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

func (s *Server) removeClient(clientID string) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.connections[clientID]
	if ok {
		delete(s.connections, clientID)
	}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request, connectionID string) {
	s.addClient(connectionID, rw, req)
	defer func() {
		s.removeClient(req.RemoteAddr)
	}()

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("client connected ID:", connectionID)

	<-req.Context().Done()
}

func (s *Server) Send(event string, message string) {
	s.RLock()
	defer s.RUnlock()

	msgBytes := []byte(fmt.Sprintf("event: %s\n\ndata:%s\n\n", event, message))
	for cID, connection := range s.connections {
		_, err := connection.writer.Write(msgBytes)
		if err != nil {
			s.removeClient(cID)
			continue
		}

		connection.flusher.Flush()
	}
}
