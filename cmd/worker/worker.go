package main

import (
	"log"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/worker"
)

func main() {
	// load configs
	cfg, err := config.InitConfig()
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

	// get rabbitmq client
	rabbitMqQueue := service.NewTaskQueue(cfg.RabbitMqDsn, cfg.RabbitMqQueueName)

	// face match and ocr service
	faceMatchService := worker.NewFaceMatchService()
	ocrService := worker.NewOCRService()

	// start the worker and process the messages
	worker := worker.NewWorker(rabbitMqQueue, psqlStore, minioStore, faceMatchService, ocrService)
	worker.ProcessMessages()
}
