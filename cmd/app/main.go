package main

import (
	"log"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/server"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
)

func main() {
	// load configs
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error while config init: %v", err)
	}

	// get psql store
	psqlStore := service.NewPsqlStore(cfg.DbDsn)

	minioConn := &db.MinioConn{
		Endpoint: cfg.MinioEndpoint,
		User:     cfg.MinioUser,
		Password: cfg.MinioPassword,
		Ssl:      cfg.MinioSSL,
	}

	// get new minio storage
	minioStorage := db.NewMinioClient(minioConn)

	// craft the server address using env vars
	host := cfg.Host
	port := cfg.Port
	addr := host + ":" + port

	// init and start the server
	server := server.NewServer(addr, psqlStore, minioStorage)
	server.Run()
}
