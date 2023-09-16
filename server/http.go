package server

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func Mount(rGroup *gin.RouterGroup) {

	message := make(chan string)

	// go func() {
	// 	sendExampleMessage(message)
	// }()

	rGroup.POST("/send-event", func(ctx *gin.Context) {
		sendMessage(message, "test message")
	})

	rGroup.GET("/events", SSEHeadersMiddleware(), handleStream(message))
}

func sendMessage(messageChn chan string, message string) {
	messageChn <- message
	// fmt.Println("start send message")
	// for i := 0; i < 10; i++ {
	// 	message <- fmt.Sprintf("test message: %d", i)
	// 	time.Sleep(time.Second * 1)
	// }
}

func handleStream(message chan string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			fmt.Println("Client closed")
		}()
		// TODO:

		fmt.Println("Client connected")
		// ctx.Stream(func(w io.Writer) bool {
		// 	// Stream message to client from message channel
		// 	if msg, ok := <-message; ok {
		// 		ctx.SSEvent("message", msg)
		// 		return true
		// 	}
		// 	return false
		// })

		ctx.Stream(func(w io.Writer) bool {
			select {
			case <-ctx.Request.Context().Done():
				fmt.Printf("Client connection closed\n")
				return false
			case msg := <-message:
				ctx.SSEvent("message", msg)
				return true
			}
		})
	}
}
