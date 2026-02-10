package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Config is the configuration of the application.
type Config struct {
	Application Application `yaml:"application" validate:"required"`

	Logger Logger `yaml:"logger" validate:"required"`

	Controller Controller `yaml:"controller" validate:"required"`

	Infrastructure Infrastructure `yaml:"infrastructure" validate:"required"`
}

// LoadConfig loads the configuration from yaml or ENV.
func LoadConfig(path string) (*Config, error) {
	confFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(confFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file: %w", err)
	}

	if err = validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
