package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/surasithaof/sse/ginserver"
	"github.com/surasithaof/sse/simpleserver"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "Ok",
		})
	})

	ginSSEServer := ginserver.NewServer()
	mountGinHandler(&router.RouterGroup, ginSSEServer)

	simpleSSEServer := simpleserver.NewServer()
	mountHTTPHandler(&router.RouterGroup, simpleSSEServer)

	err := router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}

func mountHTTPHandler(rGroup *gin.RouterGroup, simpleSSEServer simpleserver.SSEServer) {
	rGroup.GET("/simple-event", func(ctx *gin.Context) {
		cID := xid.New().String()
		simpleSSEServer.ServeHTTP(ctx.Writer, ctx.Request, cID)
	})

	rGroup.POST("/simple-event", func(ctx *gin.Context) {
		simpleSSEServer.Broadcast("event", map[string]any{
			"message": "test broadcast message",
		})
		ctx.Status(200)
	})

	rGroup.POST("/simple-event/:connectionID", func(ctx *gin.Context) {
		connectionID := ctx.Param("connectionID")
		simpleSSEServer.SendMessage(connectionID, "event", map[string]any{
			"message": "test send message",
		})
		ctx.Status(200)
	})

}

func mountGinHandler(rGroup *gin.RouterGroup, ginSSEServer ginserver.SSEServer) {
	rGroup.GET("/gin-events", ginserver.SSEHeadersMiddleware(), func(ctx *gin.Context) {
		connectionID := xid.New().String()
		ginSSEServer.Listen(ctx, connectionID)
	})

	rGroup.POST("/gin-event", func(ctx *gin.Context) {
		ginSSEServer.BroadcastMessage(ginserver.Event{
			Event: "event",
			Message: map[string]any{
				"message": "test broadcast message",
			},
		})
		ctx.Status(200)
	})

	rGroup.POST("/gin-event/:clientID", func(ctx *gin.Context) {
		clientID := ctx.Param("clientID")
		err := ginSSEServer.SendMessage(clientID, ginserver.Event{
			Event: "event",
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
