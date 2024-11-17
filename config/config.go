package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host          string
	Port          string
	DbDsn         string
	HashPassword  string
	MinioUser     string
	MinioPassword string
	MinioEndpoint string
	MinioSSL      string
	MinioBucket   string
}

var Envs = initConfig()

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		Host:          getEnv("HOST", ""),
		Port:          getEnv("PORT", "8080"),
		DbDsn:         getEnv("DB_DSN", ""),
		HashPassword:  getEnv("HASH_PASSWORD", ""),
		MinioUser:     getEnv("MINIO_USER", ""),
		MinioPassword: getEnv("MINIO_PASSWORD", ""),
		MinioEndpoint: getEnv("MINIO_ENDPOINT", ""),
		MinioSSL:      getEnv("MINIO_SSL", ""),
		MinioBucket:   getEnv("MINIO_BUCKET_NAME", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
