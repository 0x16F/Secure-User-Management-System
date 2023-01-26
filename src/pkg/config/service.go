package config

import (
	"os"

	"github.com/goccy/go-json"
)

func NewConfig() (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile("configs/config.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}
