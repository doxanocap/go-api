package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func InitConfig() (*Config, error) {
	if err := LoadFile(".env"); err != nil {
		log.Println(fmt.Errorf("config: %w", err))
	}

	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("envconfig.Process: %w", err)
	}

	return &cfg, nil
}
