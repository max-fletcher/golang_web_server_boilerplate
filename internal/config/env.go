package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port             string
	DbURL            string
	AppMode          string
	LocalBaseUrl     string
	LiveBaseUrl      string
	RedisURL         string
	RedisCacheExpiry time.Duration
	CacheActive      bool
	JWTSecret        string
	JWTExpiry        time.Duration
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

	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	JWT_EXPIRY := os.Getenv("JWT_EXPIRY")
	if JWT_EXPIRY == "" {
		return nil, fmt.Errorf("JWT_EXPIRY is not set")
	}
	JWT_EXPIRY_INT, err := strconv.Atoi(JWT_EXPIRY)
	if err != nil {
		return nil, fmt.Errorf("Error converting JWT_EXPIRY to int")
	}
	JWT_EXPIRY_DURATION := time.Duration(JWT_EXPIRY_INT) * time.Minute

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

	REDIS_URL := os.Getenv("REDIS_URL")
	if REDIS_URL == "" {
		return nil, fmt.Errorf("REDIS_URL is not set")
	}

	REDIS_CACHE_EXPIRY_STRING := os.Getenv("REDIS_CACHE_EXPIRY")
	if REDIS_CACHE_EXPIRY_STRING == "" {
		return nil, fmt.Errorf("REDIS_CACHE_EXPIRY is not set")
	}
	REDIS_CACHE_EXPIRY_MINUTES, err := time.ParseDuration(REDIS_CACHE_EXPIRY_STRING + "m")
	if err != nil {
		return nil, fmt.Errorf("REDIS_CACHE_EXPIRY is not a valid value")
	}

	CACHE_ACTIVE_STRING := os.Getenv("CACHE_ACTIVE")
	if CACHE_ACTIVE_STRING == "" {
		return nil, fmt.Errorf("CACHE_ACTIVE is not set")
	}
	var CACHE_ACTIVE_BOOL bool
	if CACHE_ACTIVE_STRING == "true" {
		CACHE_ACTIVE_BOOL = true
	} else if CACHE_ACTIVE_STRING == "false" {
		CACHE_ACTIVE_BOOL = false
	} else {
		return nil, fmt.Errorf("REDIS_CACHE_EXPIRY is not a valid value")
	}

	return &Config{
		Port:             port,
		DbURL:            DB_URL,
		AppMode:          APP_MODE,
		LocalBaseUrl:     LOCAL_BASE_URL,
		LiveBaseUrl:      LIVE_BASE_URL,
		RedisURL:         REDIS_URL,
		RedisCacheExpiry: REDIS_CACHE_EXPIRY_MINUTES,
		CacheActive:      CACHE_ACTIVE_BOOL,
		JWTSecret:        JWT_SECRET,
		JWTExpiry:        JWT_EXPIRY_DURATION,
	}, nil
}
