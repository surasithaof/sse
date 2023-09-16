package simplehttp

import "net/http"

type SSEServer interface {
	ServeHTTP(rw http.ResponseWriter, req *http.Request, connectionID string)
	SendMessage(connectionID string, event string, message any)
	Broadcast(event string, message any)
}
