package config

import (
	"github.com/karthiknatarajan/plinth/internal/types"
	"github.com/kelseyhightower/envconfig"
)

func LoadConfig() (*types.Config, error) {
	config := new(types.Config)

	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}

	config.Process()

	return config, nil
}
