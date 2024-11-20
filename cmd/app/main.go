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

	// craft the server address using env vars
	host := cfg.Host
	port := cfg.Port
	addr := host + ":" + port

	// init and start the server
	server := server.NewServer(addr, psqlStore, minioStore, redisStore)
	server.Run()
}
