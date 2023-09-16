package gin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	sseGin "github.com/surasithaof/sse/gin"
)

const (
	EventContextPath = "/events"
	TestConnectionID = "testConnectionID"
)

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func (r *TestResponseRecorder) closeClient() {
	r.closeChannel <- true
}

func CreateTestResponseRecorder() *TestResponseRecorder {
	return &TestResponseRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func TestSSEServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	res := CreateTestResponseRecorder()
	_, r := gin.CreateTestContext(res)

	sseServer := sseGin.NewServer()

	r.GET(EventContextPath, func(ctx *gin.Context) {
		sseServer.Listen(ctx, TestConnectionID)
	})

	srv := httptest.NewServer(r)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL+EventContextPath, nil)
	if err != nil {
		t.Error(err)
	}

	messageChan := make(chan string)
	defer close(messageChan)

	go func() {
		client := &http.Client{}
		client.Timeout = 5 * time.Second
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		// Read the SSE message from the response
		reader := resp.Body
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil {
			t.Error(err)
		}

		message := string(buf[:n])
		messageChan <- message
	}()

	// wait for client connect
	time.Sleep(2 * time.Second)

	err = sseServer.SendMessage(TestConnectionID, sseGin.Event{
		Event:   "message",
		Message: "test message",
	})
	if err != nil {
		t.Error(err)
	}

	tc := time.NewTicker(5 * time.Second)

	expected := "event:message\ndata:test message\n\n"
	select {
	case msg := <-messageChan:
		fmt.Println(msg)
		if msg != expected {
			t.Errorf("Expected SSE message: %s, but got: %s", expected, msg)
		}
	case <-tc.C:
		t.Error("time out")
	}
}
