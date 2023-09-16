package main

import (
	"fmt"
	"io"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	simplehttp "github.com/surasithaof/sse/http"
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

	simpleHTTPServer := simplehttp.NewServer()

	router.GET("/simple-event", func(ctx *gin.Context) {
		cID := xid.New().String()
		simpleHTTPServer.ServeHTTP(ctx.Writer, ctx.Request, cID)
	})

	router.POST("/simple-event", func(ctx *gin.Context) {
		simpleHTTPServer.Broadcast("event", gin.H{
			"message": "test broadcast message",
		})
	})

	router.POST("/simple-event/:connectionID", func(ctx *gin.Context) {
		connectionID := ctx.Param("connectionID")
		simpleHTTPServer.SendMessage(connectionID, "event", gin.H{
			"message": "test send message",
		})
	})

	err := router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}

func mountHandler(rGroup *gin.RouterGroup, sseServer server.SSEServer) {
	rGroup.POST("/send-event", func(ctx *gin.Context) {
		sseServer.BroadcastMessage(server.Event{
			Event: "event",
			Message: gin.H{
				"message": "test broadcast message",
			},
		})
		ctx.Status(200)
	})

	rGroup.POST("/send-event/:clientID", func(ctx *gin.Context) {
		clientID := ctx.Param("clientID")
		client, found := sseServer.Connection(clientID)
		if !found {
			ctx.AbortWithStatusJSON(404, gin.H{
				"error": "not_found_client", // TODO: need to handle error
			})
			return
		}

		sseServer.SendMessage(client.ID, server.Event{
			Event:   "event",
			Message: "test send message",
		})
		ctx.Status(200)
	})

	rGroup.GET("/events", server.SSEHeadersMiddleware(), sseHandler(sseServer))
	// rGroup.GET("/events", func(ctx *gin.Context) {
	// 	sseServer.ServeHTTP(ctx.Writer, ctx.Request)
	// })
}

func sseHandler(sseServer server.SSEServer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		connectionID := xid.New().String()
		client := sseServer.AddConnection(connectionID)
		fmt.Println("Client connected ID:", client.ID)

		defer func() {
			sseServer.RemoveConnection(connectionID)
			fmt.Println("Client disconnected")
		}()

		ctx.Stream(func(w io.Writer) bool {
			select {
			case <-ctx.Request.Context().Done():
				fmt.Printf("Client connection closed\n")
				return false
			case event := <-client.EventChan():
				ctx.SSEvent(event.Event, event.Message)
				return true
			}
		})
	}
}
