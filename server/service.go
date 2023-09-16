package server

import (
	"sync"
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

func (s *serverImpl) AddConnection(connectionID string) Connection {
	s.Lock()
	defer s.Unlock()

	connection := Connection{
		ID:          connectionID,
		messageChan: make(chan Event),
	}

	s.connections[connection.ID] = &connection
	return connection
}

func (s *serverImpl) RemoveConnection(connectionID string) bool {
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

func (s *serverImpl) Connections() map[string]*Connection {
	return s.connections
}

func (s *serverImpl) Connection(connectionID string) (*Connection, bool) {
	connection, ok := s.connections[connectionID]
	if !ok {
		return nil, false
	}
	return connection, true
}

func (s *serverImpl) SendMessage(connectionID string, event Event) {
	s.Lock()
	defer s.Unlock()

	connection, ok := s.connections[connectionID]
	if !ok {
		return
	}
	connection.SendMessage(event)
}

func (s *serverImpl) BroadcastMessage(event Event) {
	s.Lock()
	defer s.Unlock()

	for _, connection := range s.connections {
		connection.SendMessage(event)
	}
}
