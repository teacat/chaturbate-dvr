package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/teacat/chaturbate-dvr/domain/stream"
	"github.com/teacat/chaturbate-dvr/gateway/handler"
)

type Gateway struct {
	routes map[string]gin.HandlerFunc
}

func New(manager stream.Manager) *Gateway {
	g := &Gateway{
		routes: make(map[string]gin.HandlerFunc),
	}
	g.routes = map[string]gin.HandlerFunc{
		"/api/list_streams":  handle(handler.NewListStreamsHandler(manager)),
		"/api/start_stream":  handle(handler.NewStartStreamHandler(manager)),
		"/api/stop_stream":   handle(handler.NewStopStreamHandler(manager)),
		"/api/pause_stream":  handle(handler.NewPauseStreamHandler(manager)),
		"/api/resume_stream": handle(handler.NewResumeStreamHandler(manager)),
	}
	return g
}

func (g *Gateway) Routes() map[string]gin.HandlerFunc {
	return g.routes
}
