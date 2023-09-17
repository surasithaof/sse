package gin

import "github.com/surasithaof/sse/shared"

type Connection struct {
	ID          string
	messageChan chan shared.Event
}

func (c *Connection) EventChan() chan shared.Event {
	return c.messageChan
}

func (c *Connection) SendMessage(event shared.Event) {
	c.messageChan <- event
}
