package config

import (
	"github.com/go-sphere/confstore"
	"github.com/go-sphere/confstore/codec"
	"github.com/go-sphere/confstore/provider"
	"github.com/go-sphere/confstore/provider/file"
	"github.com/go-sphere/confstore/provider/http"
	"github.com/go-sphere/sphere-simple-layout/internal/server/api"
	"github.com/go-sphere/sphere/log/zapx"
	"github.com/go-sphere/sphere/utils/secure"
)

var BuildVersion = "dev"

type Config struct {
	Environments map[string]string `json:"environments" yaml:"environments"`
	Log          zapx.Config       `json:"log" yaml:"log"`
	API          api.Config        `json:"api" yaml:"api"`
}

func NewEmptyConfig() Config {
	return Config{
		Environments: map[string]string{},
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
			JWT: secure.RandString(32),
			HTTP: api.HTTPConfig{
				Address: "0.0.0.0:8899",
			},
		},
	}
}

func NewConfig(path string) (*Config, error) {
	config, err := confstore.Load[Config](provider.NewSelect(
		path,
		provider.If(file.IsLocalPath, func(s string) provider.Provider {
			return file.New(path, file.WithExpandEnv())
		}),
		provider.If(http.IsRemoteURL, func(s string) provider.Provider {
			return http.New(path, http.WithTimeout(10))
		}),
	), codec.JsonCodec())
	if err != nil {
		return nil, err
	}
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	return config, nil
}
