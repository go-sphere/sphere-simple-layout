package config

import (
	"github.com/TBXark/confstore"
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

func NewConfig(path string) (*Config, error) {
	config, err := confstore.Load[Config](path)
	if err != nil {
		return nil, err
	}
	return setDefaultConfig(config), nil
}
