package gin

import "github.com/gin-gonic/gin"

type SSEServer interface {
	SendMessage(connectionID string, event Event) error
	BroadcastMessage(event Event)
	Listen(ctx *gin.Context, connectionID string)
}
