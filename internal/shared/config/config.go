package config

import "os"

type AppConfig struct {
	MongoURL    string
	Database    string
	Port        string
	AdminAPIKey string
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

	return AppConfig{
		MongoURL:    os.Getenv("MONGO_URL"),
		Database:    database,
		Port:        port,
		AdminAPIKey: os.Getenv("ADMIN_API_KEY"),
	}
}
