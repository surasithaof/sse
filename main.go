package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	server.MountHandler(&router.RouterGroup, sseServer)

	err := router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
