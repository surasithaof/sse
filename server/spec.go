package server

type SSEServer interface {
	AddClient(clientID string) Client
	RemoveClient(clientID string) bool
	Clients() map[string]*Client
	Client(clientID string) (*Client, bool)

	SendMessage(clientID string, message string)
	BroadcastMessage(message string)
}
