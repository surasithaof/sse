package server

type Connection struct {
	ID          string
	messageChan chan Event
}

type Event struct {
	Event   string
	Message any
}

func (c *Connection) EventChan() chan Event {
	return c.messageChan
}

func (c *Connection) SendMessage(event Event) {
	c.messageChan <- event
}
