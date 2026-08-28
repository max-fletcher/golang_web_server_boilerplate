package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port  string
	DbURL string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT is not set")
	}

	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		return nil, fmt.Errorf("DB_URL is not set")
	}

	return &Config{
		Port:  port,
		DbURL: DB_URL,
	}, nil
}
