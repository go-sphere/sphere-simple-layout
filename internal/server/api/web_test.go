package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	apiv1 "github.com/go-sphere/sphere-simple-layout/api/api/v1"
	service "github.com/go-sphere/sphere-simple-layout/internal/service/api"
)

// startTestWeb assembles and starts the server exactly as the app does —
// NewWebServer then Start on a real loopback listener — and returns its base
// URL. Start registers the routes and then blocks in ListenAndServe, so the
// test waits for the port to accept connections instead of touching the engine
// while Start is still mounting routes on it.
func startTestWeb(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	web := NewWebServer(Config{HTTP: HTTPConfig{Address: address}}, service.NewService())
	startErr := make(chan error, 1)
	go func() { startErr <- web.Start(context.Background()) }()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = web.Stop(ctx)
		select {
		case err := <-startErr:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				t.Errorf("server start: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Error("server did not stop")
		}
	})

	deadline := time.Now().Add(2 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return "http://" + address
		}
		select {
		case err := <-startErr:
			t.Fatalf("server start: %v", err)
		default:
		}
		if time.Now().After(deadline) {
			t.Fatalf("server did not accept connections: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type envelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

var client = &http.Client{Timeout: 2 * time.Second}

func do(t *testing.T, baseURL, method, path, body string) (int, envelope) {
	t.Helper()
	req, err := http.NewRequest(method, baseURL+path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("%s %s: body is not JSON: %v; body=%q", method, path, err, raw)
	}
	return resp.StatusCode, env
}

// TestDeclaredRoutesAreMounted pins that every endpoint the generated code
// declares is reachable on the assembled server: an unregistered route would
// answer 404 or 405 rather than reach the handler.
func TestDeclaredRoutesAreMounted(t *testing.T) {
	baseURL := startTestWeb(t)
	for _, endpoint := range apiv1.EndpointsGreetService {
		method, path := endpoint[1], endpoint[2]
		status, _ := do(t, baseURL, method, path, `{}`)
		if status == http.StatusNotFound || status == http.StatusMethodNotAllowed {
			t.Errorf("%s %s (%s) status = %d, want the route to be mounted", method, path, endpoint[0], status)
		}
	}
}

func TestGreetRouteContract(t *testing.T) {
	baseURL := startTestWeb(t)

	t.Run("valid request", func(t *testing.T) {
		status, env := do(t, baseURL, http.MethodPost, "/v1/greet?title=Dr", `{"name":"Ada"}`)
		if status != http.StatusOK || !env.Success {
			t.Fatalf("status = %d success = %v, want 200 true; message=%q", status, env.Success, env.Message)
		}
		var data apiv1.GreetResponse
		if err := json.Unmarshal(env.Data, &data); err != nil {
			t.Fatalf("decode data: %v", err)
		}
		if got, want := data.Message, "Hello Dr Ada!"; got != want {
			t.Errorf("message = %q, want %q", got, want)
		}
	})

	t.Run("validation failure is a 400 with the violation", func(t *testing.T) {
		status, env := do(t, baseURL, http.MethodPost, "/v1/greet", `{"name":""}`)
		if status != http.StatusBadRequest || env.Success {
			t.Fatalf("status = %d success = %v, want 400 false", status, env.Success)
		}
		if !strings.Contains(env.Message, "at least 1") {
			t.Errorf("message = %q, want the min_len violation", env.Message)
		}
	})

	t.Run("wrong method is a JSON 405", func(t *testing.T) {
		if status, _ := do(t, baseURL, http.MethodGet, "/v1/greet", ""); status != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want 405", status)
		}
	})

	t.Run("unknown path is a JSON 404", func(t *testing.T) {
		if status, _ := do(t, baseURL, http.MethodPost, "/v1/missing", `{}`); status != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", status)
		}
	})
}
