package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	DbURL        string
	AppMode      string
	LocalBaseUrl string
	LiveBaseUrl  string
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

	APP_MODE := os.Getenv("APP_MODE")
	if APP_MODE == "" {
		return nil, fmt.Errorf("APP_MODE is not set")
	}

	LOCAL_BASE_URL := os.Getenv("LOCAL_BASE_URL")
	if LOCAL_BASE_URL == "" {
		return nil, fmt.Errorf("LOCAL_BASE_URL is not set")
	}

	LIVE_BASE_URL := os.Getenv("LIVE_BASE_URL")
	if LIVE_BASE_URL == "" {
		return nil, fmt.Errorf("LIVE_BASE_URL is not set")
	}

	return &Config{
		Port:         port,
		DbURL:        DB_URL,
		AppMode:      APP_MODE,
		LocalBaseUrl: LOCAL_BASE_URL,
		LiveBaseUrl:  LIVE_BASE_URL,
	}, nil
}
