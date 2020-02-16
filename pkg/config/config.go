package config

import (
	"io"
	"net/url"

	"github.com/caarlos0/env/v6"
)

// Config stores configuration parameters for interacting with the server at a
// global level. This can include listening address, feature flags and other
// configurations.
type Config struct {
	Addr       string  `env:"ADDR" envDefault:"0.0.0.0:8080"`
	ConfigPath string  `env:"CONFIG_PATH"`
	ConfigURL  url.URL `env:"CONFIG_URL"`
}

// New initializes a Config, attempting to parse parames from Envs.
func New() (Config, error) {
	c := Config{}

	if err := env.Parse(&c); err != nil {
		return c, err
	}

	return c, nil
}

// Load attempts to fetch a Router configuration from one of the optional
// locations (URL or Filepath). On success it returns an io.Reader for this
// file otherwise an error is returned.
func (c *Config) Load() (io.Reader, error) {
	return nil, nil
}
