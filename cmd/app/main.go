package main

import (
	"log"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	_ "github.com/justsushant/one2n-go-bootcamp/go-ekyc/docs"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/server"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
)

// @title           Ekyc REST API
// @version         1.0
// @basePath 		/api/v1
func main() {
	// load configs
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Error while config init: %v", err)
	}

	// get psql store
	psqlStore := service.NewPsqlStore(cfg.DbDsn)

	// get minio store
	minioConn := &db.MinioConn{
		Endpoint: cfg.MinioEndpoint,
		User:     cfg.MinioUser,
		Password: cfg.MinioPassword,
		Ssl:      cfg.MinioSSL,
	}
	minioStore := service.NewMinioStore(minioConn, cfg.MinioBucket)

	// get redis stores
	redisStore := service.NewRedisStore(cfg.RedisDsn)

	// get rabbitmq client
	rabbitMqQueue := service.NewTaskQueue(cfg.RabbitMqDsn, cfg.RabbitMqQueueName)

	// craft the server address using env vars
	host := cfg.Host
	port := cfg.Port
	addr := host + ":" + port

	// init and start the server
	server := server.New(&server.ServerConfig{
		Addr:       addr,
		DataStore:  psqlStore,
		FileStore:  minioStore,
		CacheStore: redisStore,
		Queue:      rabbitMqQueue,
	})
	server.Run()
}
