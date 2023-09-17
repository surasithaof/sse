package http

import (
	"net/http"

	"github.com/surasithaof/sse/shared"
)

//go:generate mockgen -source=./spec.go -destination=./mocks/server.go -package=mocks "github.com/surasithaof/sse" SSEServer

type SSEServer interface {
	Listen(rw http.ResponseWriter, req *http.Request, connectionID string)
	SendMessage(connectionID string, event shared.Event) error
	Broadcast(event shared.Event)
}
