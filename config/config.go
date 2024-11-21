package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Host              string `env:"HOST"`
	Port              string `env:"PORT,required"`
	DbDsn             string `env:"DB_DSN,required"`
	HashPassword      string `env:"HASH_PASSWORD,required"`
	MinioUser         string `env:"MINIO_USER,required"`
	MinioPassword     string `env:"MINIO_PASSWORD,required"`
	MinioEndpoint     string `env:"MINIO_ENDPOINT,required"`
	MinioSSL          bool   `env:"MINIO_SSL,required"`
	MinioBucket       string `env:"MINIO_BUCKET_NAME,required"`
	RedisDsn          string `env:"REDIS_DSN,required"`
	RabbitMqDsn       string `env:"RABBITMQ_DSN,required"`
	RabbitMqQueueName string `env:"RABBITMQ_QUEUE_NAME,required"`
}

func InitConfig() (*Config, error) {
	// load the env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// parse the env file
	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &cfg, nil
}
