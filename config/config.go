package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigDirName  = "gh-ci"
	ConfigFileName = "config"
	ConfigFileExt  = "yaml"
)

type Config struct {
	Github GithubConfig
}

type GithubConfig struct {
	Repositories []string
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error finding home directory: %w", err)
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(home, ".config")
	}
	configDir := filepath.Join(configHome, ConfigDirName)

	viper.AddConfigPath(configDir)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileExt)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(configDir); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			return nil, fmt.Errorf("default config created at %s", configDir)
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Github.Repositories) == 0 {
		return fmt.Errorf("no repositories found in config")
	}
	for _, repo := range c.Github.Repositories {
		if len(repo) == 0 {
			return fmt.Errorf("repository name cannot be empty")
		}
		parts := strings.Split(repo, "/")
		if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
			return fmt.Errorf("repository name must be in the format 'owner/repo'")
		}
	}
	return nil
}

func createDefaultConfig(configDir string) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.%s", ConfigFileName, ConfigFileExt))

	yaml, _ := yaml.Marshal(Config{})
	newConfigFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() {
		closeErr := newConfigFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	_, err = newConfigFile.WriteString(string(yaml))
	if err != nil {
		return fmt.Errorf("failed to write to config file: %w", err)
	}

	fmt.Printf("Created default configuration at: %s\n", configPath)
	fmt.Println("Please update it with your repository names.")
	fmt.Println("Example: \n  - owner/repo1\n  - owner/repo2")

	return nil
}
