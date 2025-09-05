package config

import (
	"context"
	"errors"

	"github.com/go-sphere/confstore"
	"github.com/go-sphere/confstore/codec"
	"github.com/go-sphere/confstore/provider"
	"github.com/go-sphere/sphere-simple-layout/internal/server/api"
	"github.com/go-sphere/sphere/log"
	"github.com/go-sphere/sphere/utils/secure"
)

var BuildVersion = "dev"

type Config struct {
	Environments map[string]string `json:"environments" yaml:"environments"`
	Log          *log.Config       `json:"log" yaml:"log"`
	API          *api.Config       `json:"api" yaml:"api"`
}

func NewEmptyConfig() *Config {
	return &Config{
		Environments: map[string]string{},
		Log: &log.Config{
			File: &log.FileConfig{
				FileName:   "./var/log/sphere.log",
				MaxSize:    10,
				MaxBackups: 10,
				MaxAge:     10,
			},
			Console: &log.ConsoleConfig{},
			Level:   "info",
		},
		API: &api.Config{
			JWT: secure.RandString(32),
			HTTP: api.HTTPConfig{
				Address: "0.0.0.0:8899",
			},
		},
	}
}

func setDefaultConfig(config *Config) *Config {
	if config.Log == nil {
		config.Log = log.NewDefaultConfig()
	}
	return config
}

func newConfProvider(path string) (provider.Provider, error) {
	if provider.IsRemoteURL(path) {
		return provider.NewHTTP(path, provider.WithTimeout(10)), nil
	}
	if provider.IsLocalPath(path) {
		return provider.NewFile(path, provider.WithExpandEnv()), nil
	}
	return nil, errors.New("unsupported config path")
}

func NewConfig(path string) (*Config, error) {
	pro, err := newConfProvider(path)
	if err != nil {
		return nil, err
	}
	config, err := confstore.Load[Config](context.Background(), pro, codec.JsonCodec())
	if err != nil {
		return nil, err
	}
	return setDefaultConfig(config), nil
}
