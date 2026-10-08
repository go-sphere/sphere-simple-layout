package config

import (
	"fmt"

	"github.com/go-sphere/confstore"
	"github.com/go-sphere/confstore/codec"
	"github.com/go-sphere/confstore/provider/file"
	"github.com/go-sphere/sphere-simple-layout/internal/server/api"
	"github.com/go-sphere/sphere/log/zapx"
)

var BuildVersion = "dev"

type Config struct {
	Log zapx.Config `json:"log" yaml:"log"`
	API api.Config  `json:"api" yaml:"api"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Log: zapx.Config{
			File: zapx.FileConfig{
				FileName:   "./var/log/sphere.log",
				MaxSize:    10,
				MaxBackups: 10,
				MaxAge:     10,
			},
			Console: zapx.ConsoleConfig{},
			Level:   "info",
		},
		API: api.Config{
			HTTP: api.HTTPConfig{
				Address: "0.0.0.0:8899",
			},
		},
	}
}

func NewConfig(path string) (*Config, error) {
	config, err := confstore.Load[Config](file.New(path), codec.JsonCodec())
	if err != nil {
		return nil, err
	}
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if err := config.API.HTTP.Validate(); err != nil {
		return nil, fmt.Errorf("api http: %w", err)
	}
	return config, nil
}
