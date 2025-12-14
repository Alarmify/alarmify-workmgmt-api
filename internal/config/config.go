package config

import (
	"os"
)

type Config struct {
	Port        string
	Environment string
	ReadTimeout int
	WriteTimeout int
	Database    DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8089"),
		Environment: getEnv("ENVIRONMENT", "development"),
		ReadTimeout: 15,
		WriteTimeout: 15,
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "workmgmt-api_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
