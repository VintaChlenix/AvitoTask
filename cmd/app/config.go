package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		URL string `yaml:"SERVER_URL"`
	} `yaml:"server"`
	Database struct {
		ConnectionString string `yaml:"DATABASE_URL"`
	} `yaml:"database"`
}

func GetConfig() (*Config, error) {
	f, err := os.Open("configs/config.yml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
