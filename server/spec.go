package server

type SSEServer interface {
	AddConnection(connectionID string) Connection
	RemoveConnection(connectionID string) bool
	Connections() map[string]*Connection
	Connection(connectionID string) (*Connection, bool)

	SendMessage(connectionID string, event Event)
	BroadcastMessage(event Event)
}
