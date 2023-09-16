package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	sseGin "github.com/surasithaof/sse/gin"
	sseHTTP "github.com/surasithaof/sse/http"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.StaticFS("/client", http.Dir("./example/public"))

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "Ok",
		})
	})

	ginSSEServer := sseGin.NewServer()
	mountGinHandler(&router.RouterGroup, ginSSEServer)

	simpleSSEServer := sseHTTP.NewServer()
	mountHTTPHandler(&router.RouterGroup, simpleSSEServer)

	err := router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}

func mountHTTPHandler(rGroup *gin.RouterGroup, simpleSSEServer sseHTTP.SSEServer) {
	rGroup.GET("/simple-events", func(ctx *gin.Context) {
		connectionID := xid.New().String()
		fmt.Println("client connection ID:", connectionID)
		simpleSSEServer.ServeHTTP(ctx.Writer, ctx.Request, connectionID)
	})

	rGroup.POST("/simple-events", func(ctx *gin.Context) {
		simpleSSEServer.Broadcast("message", map[string]any{
			"message": "test broadcast message",
		})
		ctx.Status(200)
	})

	rGroup.POST("/simple-events/:connectionID", func(ctx *gin.Context) {
		connectionID := ctx.Param("connectionID")
		err := simpleSSEServer.SendMessage(connectionID, "message", map[string]any{
			"message": "test send message",
		})
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"error": err.Error(),
			})
		}
		ctx.Status(200)
	})

}

func mountGinHandler(rGroup *gin.RouterGroup, ginSSEServer sseGin.SSEServer) {
	rGroup.GET("/gin-events", sseGin.SSEHeadersMiddleware(), func(ctx *gin.Context) {
		connectionID := xid.New().String()
		fmt.Println("client connection ID:", connectionID)
		ginSSEServer.Listen(ctx, connectionID)
	})

	rGroup.POST("/gin-events", func(ctx *gin.Context) {
		ginSSEServer.BroadcastMessage(sseGin.Event{
			Event: "message",
			Message: map[string]any{
				"message": "test broadcast message",
			},
		})
		ctx.Status(200)
	})

	rGroup.POST("/gin-events/:clientID", func(ctx *gin.Context) {
		clientID := ctx.Param("clientID")
		err := ginSSEServer.SendMessage(clientID, sseGin.Event{
			Event: "message",
			Message: map[string]any{
				"message": "test send message",
			},
		})

		// TODO: need to handle error
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"error": err.Error(),
			})
		}
		ctx.Status(200)
	})
}
