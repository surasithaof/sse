package gin

import (
	"errors"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/surasithaof/sse/shared"
)

type serverImpl struct {
	connections map[string]*Connection
	messageChan chan string
	sync.Mutex
}

func NewServer() SSEServer {
	s := serverImpl{
		connections: make(map[string]*Connection),
		messageChan: make(chan string),
	}
	return &s
}

func (s *serverImpl) addConnection(connectionID string) Connection {
	s.Lock()
	defer s.Unlock()

	connection := Connection{
		ID:          connectionID,
		messageChan: make(chan shared.Event),
	}

	s.connections[connection.ID] = &connection
	return connection
}

func (s *serverImpl) removeConnection(connectionID string) bool {
	s.Lock()
	defer s.Unlock()

	connection, ok := s.connections[connectionID]
	if !ok {
		return false
	}
	close(connection.messageChan)
	delete(s.connections, connection.ID)
	return true
}

func (s *serverImpl) SendMessage(connectionID string, event shared.Event) error {
	s.Lock()
	defer s.Unlock()

	connection, ok := s.connections[connectionID]
	if !ok {
		// TODO: need to handle error
		return errors.New("not_found_connection")
	}
	connection.SendMessage(event)
	return nil
}

func (s *serverImpl) BroadcastMessage(event shared.Event) {
	s.Lock()
	defer s.Unlock()

	for _, connection := range s.connections {
		connection.SendMessage(event)
	}
}

func (s *serverImpl) Listen(ctx *gin.Context, connectionID string) {
	client := s.addConnection(connectionID)

	defer func() {
		s.removeConnection(connectionID)
	}()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Request.Context().Done():
			return false
		case event := <-client.EventChan():
			ctx.SSEvent(event.Event, event.Message)
			return true
		}
	})
}
