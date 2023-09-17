package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/surasithaof/sse/shared"
)

//go:generate mockgen -source=./spec.go -destination=./mocks/server.go -package=mocks "github.com/surasithaof/sse" SSEServer

type SSEServer interface {
	SendMessage(connectionID string, event shared.Event) error
	BroadcastMessage(event shared.Event)
	Listen(ctx *gin.Context, connectionID string)
}
