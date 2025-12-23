package httpsrv

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-sphere/httpx"
	"github.com/go-sphere/httpx/ginx"
	"github.com/go-sphere/sphere/log"
)

// NewGinServer initializes and returns a new HTTP server engine configured with the specified address and middlewares.
func NewGinServer(name, addr string) httpx.Engine {
	logger := log.With(log.WithAttrs(map[string]any{"module": name}), log.DisableCaller())
	engine := gin.New()
	if zapLogger, err := log.UnwrapZapLogger(logger); err == nil {
		engine.Use(ginzap.Ginzap(zapLogger, time.RFC3339, true))
		engine.Use(ginzap.RecoveryWithZap(zapLogger, true))
	} else {
		engine.Use(gin.Recovery())
	}
	app := ginx.New(
		ginx.WithEngine(engine),
		ginx.WithServerAddr(addr),
	)
	return app
}
