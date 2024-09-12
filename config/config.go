package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App          `yaml:"app"`
		HTTP         `yaml:"http"`
		GRPC         `yaml:"grpc"`
		Orchestrator `yaml:"orchestrator"`
	}

	App struct {
		Name     string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version  string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		LogLevel string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}
	HTTP struct {
		Port int32 `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	GRPC struct {
		Port int32 `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	Orchestrator struct {
		Interval float32 `env-required:"true" yaml:"interval" env:"ORCHESTRATOR_INTERVAL"`
	}
)

// NewConfig returns app config.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	if path == "" {
		fmt.Println("Config path not provided")
		return nil, fmt.Errorf("config path not provided")
	}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
