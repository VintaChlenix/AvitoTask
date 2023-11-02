package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		URL string `yaml:"SERVER_URL"`
	} `yaml:"server"`
	Database struct {
		ConnectionString string `yaml:"DATABASE_URL"`
	} `yaml:"database"`
}

func GetConfig(path string) (cfg *Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		ferr := f.Close()
		if ferr == nil {
			return
		}
		err = errors.Join(ferr, err)
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}
