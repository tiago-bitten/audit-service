package config

import (
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	MongoURL          string
	Database          string
	Port              string
	AdminAPIKey       string
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

func LoadConfig() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database := os.Getenv("MONGO_DATABASE")
	if database == "" {
		database = "auditservice"
	}

	rateLimitRequests := 100
	if v := os.Getenv("RATE_LIMIT_REQUESTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			rateLimitRequests = n
		}
	}

	rateLimitWindow := 60 * time.Second
	if v := os.Getenv("RATE_LIMIT_WINDOW"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			rateLimitWindow = d
		}
	}

	return AppConfig{
		MongoURL:          os.Getenv("MONGO_URL"),
		Database:          database,
		Port:              port,
		AdminAPIKey:       os.Getenv("ADMIN_API_KEY"),
		RateLimitRequests: rateLimitRequests,
		RateLimitWindow:   rateLimitWindow,
	}
}
