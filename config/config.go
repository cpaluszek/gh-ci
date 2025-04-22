package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Github GithubConfig
}

type GithubConfig struct {
	Token string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.Github.Token == "" {
		return nil, fmt.Errorf("github token is required. set it in config file")
	}

	return &cfg, nil
}
