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
)

func init() {
	httpz.SetDefaultErrorParser(func(err error) (int32, int32, string) {
		var ve *protovalidate.ValidationError
		if errors.As(err, &ve) {
			msgs := make([]string, 0, len(ve.Violations))
			for _, v := range ve.Violations {
				msgs = append(msgs, v.Proto.GetMessage())
			}
			return 0, http.StatusBadRequest, strings.Join(msgs, ",")
		}
		return httpx.ParseError(err)
	})
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
