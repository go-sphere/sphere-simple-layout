package main

import (
	"github.com/go-sphere/sphere-simple-layout/internal/server/api"
	"github.com/go-sphere/sphere/core/boot"
)

func newApplication(
	api *api.Web,
) *boot.Application {
	return boot.NewApplication(
		api,
	)
}
