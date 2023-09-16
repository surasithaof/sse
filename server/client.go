package server

type Client struct {
	ID          string
	messageChan chan string
}

func (c *Client) MessageChan() chan string {
	return c.messageChan
}
