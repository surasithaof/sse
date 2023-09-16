package main

import (
	"fmt"
	"io"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/surasithaof/sse/server"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "Ok",
		})
	})

	sseServer := server.NewServer()
	mountHandler(&router.RouterGroup, sseServer)

	err := router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}

func mountHandler(rGroup *gin.RouterGroup, sseServer server.SSEServer) {
	rGroup.POST("/send-event", func(ctx *gin.Context) {
		sseServer.BroadcastMessage("test broadcast message")
		ctx.Status(200)
	})

	rGroup.POST("/send-event/:clientID", func(ctx *gin.Context) {
		clientID := ctx.Param("clientID")
		client, found := sseServer.Client(clientID)
		if !found {
			ctx.AbortWithStatusJSON(404, gin.H{
				"error": "not_found_client", // TODO: need to handle error
			})
			return
		}

		sseServer.SendMessage(client.ID, "test send message")
		ctx.Status(200)
	})

	rGroup.GET("/events", server.SSEHeadersMiddleware(), sseHandler(sseServer))
}

func sseHandler(sseServer server.SSEServer) gin.HandlerFunc {
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
			case msg := <-client.MessageChan():
				ctx.SSEvent("message", msg)
				return true
			}
		})
	}
}
