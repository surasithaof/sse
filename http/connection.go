package simplehttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Connection struct {
	ID         string
	writer     http.ResponseWriter
	flusher    http.Flusher
	requestCtx context.Context
}

func (connection *Connection) send(event string, message any) error {

	msg, ok := message.(string)
	if !ok {
		msgJSON, err := json.Marshal(message)
		if err != nil {
			return err
		}
		msg = string(msgJSON)
	}

	msgBytes := []byte(fmt.Sprintf("event: %s\ndata:%s\n\n", event, msg))
	_, err := connection.writer.Write(msgBytes)
	if err != nil {
		return err
	}

	connection.flusher.Flush()
	return nil
}
