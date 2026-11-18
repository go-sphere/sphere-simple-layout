package api

import (
	"context"
	"fmt"
	"strings"

	apiv1 "github.com/go-sphere/sphere-simple-layout/api/api/v1"
)

var _ apiv1.GreetServiceHTTPServer = (*Service)(nil)

func (s *Service) Greet(ctx context.Context, request *apiv1.GreetRequest) (*apiv1.GreetResponse, error) {
	title := strings.TrimSpace(request.Title)
	if title == "" {
		return &apiv1.GreetResponse{
			Message: fmt.Sprintf("Hello %s!", request.Name),
		}, nil
	}
	return &apiv1.GreetResponse{
		Message: fmt.Sprintf("Hello %s %s!", title, request.Name),
	}, nil
}
