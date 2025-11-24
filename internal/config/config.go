package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort    string
	DatabaseURL string
}

func Load() Config {
	_ = godotenv.Load()

	httpPort := getEnv("HTTP_PORT", "8080")
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return Config{
		HTTPPort:    httpPort,
		DatabaseURL: dbURL,
	}
}

func getEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return defaultVal
}
