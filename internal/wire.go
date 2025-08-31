package internal

import (
	"github.com/go-sphere/sphere-simple-layout/internal/config"
	"github.com/go-sphere/sphere-simple-layout/internal/server"
	"github.com/go-sphere/sphere-simple-layout/internal/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	// Internal application components
	wire.NewSet(
		server.ProviderSet,
		service.ProviderSet,
		config.ProviderSet,
	),
)
