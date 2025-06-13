package main

import (
	"encoding/json"
	"os"
)

type DatabaseConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Db   string `json:"db"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
}

func loadConfiguration(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == "" {
		config.Database.Port = "5432"
	}

	if config.Database.User == "" {
		config.Database.User = "eyub"
	}

	if config.Database.Pass == "" {
		config.Database.Pass = "1234"
	}

	if config.Database.Db == "" {
		config.Database.Db = "iot_platform"
	}

	return &config, nil
}
