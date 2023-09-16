package http

import (
	"context"
	"net/http"
)

type Connection struct {
	ID         string
	writer     http.ResponseWriter
	flusher    http.Flusher
	requestCtx context.Context
}
