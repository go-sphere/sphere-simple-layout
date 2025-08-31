package api

type HTTPConfig struct {
	Address string `json:"address" yaml:"address"`
}

type Config struct {
	JWT  string     `json:"jwt" yaml:"jwt"`
	HTTP HTTPConfig `json:"http" yaml:"http"`
}
