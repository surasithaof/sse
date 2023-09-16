package server

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func MountHandler(rGroup *gin.RouterGroup, sseServer SSEServer) {
	rGroup.POST("/send-event", func(ctx *gin.Context) {
		sseServer.BroadcastMessage("test broadcast message")
		ctx.Status(200)
	})

	rGroup.POST("/send-event/:clientID", func(ctx *gin.Context) {
		cID := ctx.Param("clientID")
		client, found := sseServer.Client(cID)
		if !found {
			ctx.AbortWithStatusJSON(404, gin.H{
				"error": "not_found_client", // TODO: need to handle error
			})
			return
		}

		sseServer.SendMessage(client.ID, "test send message")
		ctx.Status(200)
	})

	rGroup.GET("/events", SSEHeadersMiddleware(), sseHandler(sseServer))
}

func sseHandler(sseServer SSEServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientID := xid.New().String()
		client := sseServer.AddClient(clientID)
		fmt.Println("Client connected ID:", client.ID)

		defer func() {
			sseServer.RemoveClient(clientID)
			fmt.Println("Client disconnected")
		}()

		ctx.Stream(func(w io.Writer) bool {
			select {
			case <-ctx.Request.Context().Done():
				fmt.Printf("Client connection closed\n")
				return false
			case msg := <-client.messageChan:
				ctx.SSEvent("message", msg)
				return true
			}
		})
	}
}
