package server

import (
	"github.com/go-sphere/sphere-simple-layout/internal/server/api"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewWebServer,
)
