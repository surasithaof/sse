package http

import "net/http"

type SSEServer interface {
	ServeHTTP(rw http.ResponseWriter, req *http.Request, connectionID string)
	SendMessage(connectionID string, event string, message any) error
	Broadcast(event string, message any)
}
