package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	appPrefix = "DPJS"
)

// Config ...
type Config struct {
	Env      string `required:"false" default:"development"`
	Port     int    `required:"false" default:"7000"`
	ProtoDir string `required:"false" default:"../protos/proto/" split_words:"true"`
}

// GetAppConfig ...
func GetAppConfig() (*Config, error) {

	var cfg Config
	var err error

	err = envconfig.Process(appPrefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
