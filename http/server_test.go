package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sseHTTP "github.com/surasithaof/sse/http"
)

const (
	EventContextPath = "/events"
	TestConnectionID = "testConnectionID"
)

func TestSSEServer(t *testing.T) {

	sseServer := sseHTTP.NewServer()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sseServer.ServeHTTP(w, r, TestConnectionID)
	}))
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

	err = sseServer.SendMessage(TestConnectionID, "message", "test message")
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
