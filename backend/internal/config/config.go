package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	ServerPort string
	JWTSecret  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
		JWTSecret:  getEnv("JWT_SECRET", "change-me"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "inventory_user"),
		DBPassword: getEnv("DB_PASSWORD", "inventory_password"),
		DBName:     getEnv("DB_NAME", "inventory_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
