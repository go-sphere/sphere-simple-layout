package service

import (
	"github.com/go-sphere/sphere-simple-layout/internal/service/api"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
)
