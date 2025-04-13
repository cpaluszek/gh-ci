package config

import (
)
// TODO: use viper

type Config struct {
	RefreshInterval int
}

func Load() (*Config, error) {
	var cfg = Config {
		RefreshInterval: 30,
	}

	return &cfg, nil
}
