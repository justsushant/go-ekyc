package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                 string
	Port                 string
	DB_DSN               string
	Access_token_secret  string
	Refresh_token_secret string
}

var Envs = initConfig()

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		Host:                 getEnv("HOST", ""),
		Port:                 getEnv("PORT", "8080"),
		DB_DSN:               getEnv("DB_DSN", ""),
		Access_token_secret:  getEnv("ACCESS_TOKEN_SECRET", ""),
		Refresh_token_secret: getEnv("REFRESH_TOKEN_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
