package api

import "github.com/go-sphere/sphere-simple-layout/internal/pkg/httpsrv"

type HTTPConfig struct {
	Address string   `json:"address" yaml:"address"`
	Cors    []string `json:"cors" yaml:"cors"`

	httpsrv.Options `yaml:",inline"`
}

type Config struct {
	HTTP HTTPConfig `json:"http" yaml:"http"`
}
