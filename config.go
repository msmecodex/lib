package utils

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWT_SECRET string
}

func loadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		JWT_SECRET: getEnv("JWT_SECRET", "Xw7S3tQvmZWZ727p7jeF7ptOCXbAfpQr"),
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
