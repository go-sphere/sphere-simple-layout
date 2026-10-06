package httpsrv_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/go-sphere/httpx"
	apiv1 "github.com/go-sphere/sphere-simple-layout/api/api/v1"
	"github.com/go-sphere/sphere-simple-layout/internal/pkg/httpsrv"
	"github.com/go-sphere/sphere/server/httpz"
)

// httpsrv installs the process-wide error parser in init(); these tests use it
// as installed. A test that swaps it must restore it and must not run in parallel.

// serveError mounts a single route that fails with err and returns the HTTP
// status and decoded error envelope a client receives for it.
func serveError(t *testing.T, err error) (int, httpz.ErrorResponse) {
	t.Helper()
	engine := httpsrv.NewServer("test", "127.0.0.1:0")
	engine.Group("/").GET("/fail", httpz.WithJson(func(httpx.Context) (string, error) {
		return "", err
	}))
	requester, ok := httpx.AsTestRequester(engine)
	if !ok {
		t.Fatalf("engine %T does not support in-process requests", engine)
	}
	resp, reqErr := requester.Do(httptest.NewRequest(http.MethodGet, "/fail", nil))
	if reqErr != nil {
		t.Fatalf("request: %v", reqErr)
	}
	defer func() { _ = resp.Body.Close() }()
	raw, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		t.Fatalf("read body: %v", readErr)
	}
	var body httpz.ErrorResponse
	if decodeErr := json.Unmarshal(raw, &body); decodeErr != nil {
		t.Fatalf("decode body: %v; body=%q", decodeErr, raw)
	}
	if body.Success {
		t.Errorf("success = true, want false; body=%q", raw)
	}
	return resp.StatusCode, body
}

func validationError(t *testing.T) error {
	t.Helper()
	// name has min_len 1, so the empty request yields one violation.
	err := protovalidate.Validate(&apiv1.GreetRequest{Name: ""})
	if _, ok := errors.AsType[*protovalidate.ValidationError](err); !ok {
		t.Fatalf("protovalidate.Validate() = %v, want a *ValidationError", err)
	}
	return err
}

func TestErrorParserRendersValidationErrorAsBadRequest(t *testing.T) {
	for name, err := range map[string]error{
		"direct":  validationError(t),
		"wrapped": fmt.Errorf("bind greet: %w", validationError(t)),
	} {
		t.Run(name, func(t *testing.T) {
			status, body := serveError(t, err)
			if status != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", status)
			}
			if body.Code != 0 {
				t.Errorf("code = %d, want 0", body.Code)
			}
			if !strings.Contains(body.Message, "at least 1") {
				t.Errorf("message = %q, want the min_len violation", body.Message)
			}
		})
	}
}

func TestErrorParserFallsBackToHTTPXParseError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantCode    int
		wantMessage string
	}{
		{
			name:        "httpx status error",
			err:         httpx.NewNotFoundError("resource missing"),
			wantStatus:  http.StatusNotFound,
			wantMessage: "resource missing",
		},
		{
			name:        "generated proto error enum",
			err:         apiv1.GreetError_GREET_ERROR_NAME_TOO_LONG,
			wantStatus:  http.StatusBadRequest,
			wantCode:    1000,
			wantMessage: "Name exceeds maximum length",
		},
		{
			// A plain error carries no status or message: it must become a 500
			// with the generic status text, never the raw err.Error().
			name:        "plain error",
			err:         errors.New("pq: password authentication failed"),
			wantStatus:  http.StatusInternalServerError,
			wantMessage: http.StatusText(http.StatusInternalServerError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body := serveError(t, tt.err)
			if status != tt.wantStatus {
				t.Fatalf("status = %d, want %d", status, tt.wantStatus)
			}
			if body.Code != tt.wantCode {
				t.Errorf("code = %d, want %d", body.Code, tt.wantCode)
			}
			if body.Message != tt.wantMessage {
				t.Errorf("message = %q, want %q", body.Message, tt.wantMessage)
			}
			if strings.Contains(body.Message, "pq:") || strings.Contains(body.Error, "pq:") {
				t.Errorf("response leaks the raw error: %+v", body)
			}
		})
	}
}
