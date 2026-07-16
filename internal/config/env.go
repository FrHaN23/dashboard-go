package config

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":5000"),

		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "3306"),
		DBUser: getEnv("DB_USERNAME", "user"),
		DBPass: getEnv("DB_PASSWORD", "password123"),
		DBName: getEnv("DB_NAME", "article"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
