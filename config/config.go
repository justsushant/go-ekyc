package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host   string
	Port   string
	DB_DSN string
}

var Envs = initConfig()

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		Host:   getEnv("HOST", ""),
		Port:   getEnv("PORT", "8080"),
		DB_DSN: getEnv("DB_DSN", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
