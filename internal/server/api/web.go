package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/go-sphere/sphere-simple-layout/api/api/v1"
	"github.com/go-sphere/sphere-simple-layout/internal/service/api"
	"github.com/go-sphere/sphere/server/ginx"
)

type Web struct {
	config  *Config
	server  *http.Server
	service *api.Service
}

func NewWebServer(conf *Config, service *api.Service) *Web {
	return &Web{
		config:  conf,
		service: service,
	}
}

func (w *Web) Identifier() string {
	return "api"
}

func (w *Web) Start(ctx context.Context) error {
	engine := gin.Default()
	apiv1.RegisterGreetServiceHTTPServer(engine, w.service)
	w.server = &http.Server{
		Addr:    w.config.HTTP.Address,
		Handler: engine.Handler(),
	}
	return ginx.Start(w.server)
}

func (w *Web) Stop(ctx context.Context) error {
	return ginx.Close(ctx, w.server)
}
