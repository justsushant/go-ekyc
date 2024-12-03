package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host              string `env:"HOST"`
	Port              string `env:"PORT,required"`
	DbDsn             string `env:"DB_DSN,required"`
	PostgresUser      string `env:"POSTGRES_USER,required"`
	PostgresPassword  string `env:"POSTGRES_PASSWORD,required"`
	PostgresEndpoint  string `env:"POSTGRES_ENDPOINT,required"`
	PostgresSSL       string `env:"POSTGRES_SSL,required"`
	PostgresDB        string `env:"POSTGRES_DB,required"`
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

func Init() (*Config, error) {
	// load the env file
	// err := godotenv.Load(".docker-compose.env")
	// if err != nil {
	// 	return nil, err
	// }

	// parse the env file
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &cfg, nil
}
