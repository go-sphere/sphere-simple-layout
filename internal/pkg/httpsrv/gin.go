package httpsrv

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"buf.build/go/protovalidate"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-sphere/httpx"
	"github.com/go-sphere/httpx/ginx"
	"github.com/go-sphere/sphere/log"
	"github.com/go-sphere/sphere/log/zapx"
	"github.com/go-sphere/sphere/server/httpz"
	"github.com/go-sphere/sphere/server/middleware/cors"
)

func init() {
	httpz.SetDefaultErrorParser(func(err error) (int32, int32, string) {
		if ve, ok := errors.AsType[*protovalidate.ValidationError](err); ok {
			msgs := make([]string, 0, len(ve.Violations))
			for _, v := range ve.Violations {
				msgs = append(msgs, v.Proto.GetMessage())
			}
			return 0, http.StatusBadRequest, strings.Join(msgs, ",")
		}
		return httpx.ParseError(err)
	})
}

// UseCORS attaches CORS middleware when origins is non-empty.
func UseCORS(engine httpx.Engine, origins []string) error {
	if len(origins) == 0 {
		return nil
	}
	mw, err := cors.NewCORS(cors.WithAllowOrigins(origins...))
	if err != nil {
		return err
	}
	engine.Use(mw)
	return nil
}

type httpxContext = httpx.Context

type jsonErrorContext struct {
	httpxContext
	gc *gin.Context
}

func (c *jsonErrorContext) JSON(code int, v any) error {
	c.gc.JSON(code, v)
	return nil
}

// NewGinServer initializes and returns a new HTTP server engine configured with the specified address and middlewares.
func NewGinServer(name, addr string) httpx.Engine {
	logger := log.With(log.WithAttrs(map[string]any{"module": name}), log.DisableCaller())
	engine := gin.New()
	if zapBackend, ok := logger.Backend().(*zapx.Backend); ok {
		engine.Use(ginzap.Ginzap(zapBackend.ZapLogger(), time.RFC3339, true))
		engine.Use(ginzap.RecoveryWithZap(zapBackend.ZapLogger(), true))
	} else {
		engine.Use(gin.Recovery())
	}
	app := ginx.New(
		ginx.WithEngine(engine),
		ginx.WithServerAddr(addr),
		ginx.WithErrorHandler(func(gc *gin.Context, err error) {
			httpz.AbortWithJsonError(&jsonErrorContext{gc: gc}, err)
		}),
	)
	return app
}
