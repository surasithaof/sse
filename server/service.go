package server

import "sync"

type serverImpl struct {
	clients     map[string]*Client
	messageChan chan string
	sync.Mutex
}

func NewServer() SSEServer {
	s := serverImpl{
		clients:     map[string]*Client{},
		messageChan: make(chan string),
	}
	return &s
}

func (s *serverImpl) AddClient(clientID string) Client {
	s.Lock()
	defer s.Unlock()

	client := Client{
		ID:          clientID,
		messageChan: make(chan string),
	}

	s.clients[client.ID] = &client
	return client
}

func (s *serverImpl) RemoveClient(clientID string) bool {
	s.Lock()
	defer s.Unlock()

	client, ok := s.clients[clientID]
	if !ok {
		return false
	}
	close(client.messageChan)
	delete(s.clients, client.ID)
	return true
}

func (s *serverImpl) Clients() map[string]*Client {
	return s.clients
}

func (s *serverImpl) Client(clientID string) (*Client, bool) {
	client, ok := s.clients[clientID]
	if !ok {
		return nil, false
	}
	return client, true
}

func (s *serverImpl) SendMessage(clientID string, message string) {
	s.Lock()
	defer s.Unlock()

	client, ok := s.clients[clientID]
	if !ok {
		return
	}
	client.messageChan <- message
}

func (s *serverImpl) BroadcastMessage(message string) {
	s.Lock()
	defer s.Unlock()

	for _, client := range s.clients {
		client.messageChan <- message
	}
}
