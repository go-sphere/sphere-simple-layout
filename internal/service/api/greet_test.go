package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apiv1 "github.com/go-sphere/sphere-simple-layout/api/api/v1"
	"github.com/go-sphere/sphere/server/ginx"
)

func TestService_Greet(t *testing.T) {
	router := gin.Default()
	apiv1.RegisterGreetServiceHTTPServer(router, &Service{})

	w := httptest.NewRecorder()
	uri := fmt.Sprintf("%s?title=%s", apiv1.EndpointsGreetService[0][2], "Mr.")
	req := httptest.NewRequest("POST", uri, bytes.NewReader([]byte(`{"name": "World"}`)))
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	var resp ginx.DataResponse[*apiv1.GreetResponse]
	err := json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedMessage := "Hello Mr. World!"
	if resp.Data.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, resp.Data.Message)
	}
}
