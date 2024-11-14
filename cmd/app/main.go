package main

import (
	"log"
	"strconv"

	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/server"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
)

func main() {
	// extracting pgsql dsn from env var
	dsn := config.Envs.DB_DSN
	if dsn == "" {
		panic("Database DSN not found")
	}

	// get new postgresql storage
	pgStorage := store.NewPostgreSQLStorage(dsn)

	// extracting minio connection vars from config (.env file)
	minioSsl, err := strconv.ParseBool(config.Envs.MinioSSL)
	if err != nil {
		log.Fatalf("minio ssl config not found")
	}
	minioConn := &store.MinioConn{
		Endpoint: config.Envs.MinioEndpoint,
		User:     config.Envs.MinioUser,
		Password: config.Envs.MinioPassword,
		Ssl:      minioSsl,
	}

	// get new minio storage
	minioStorage := store.NewMinioClient(minioConn)

	// craft the server address using env vars
	host := config.Envs.Host
	port := config.Envs.Port
	addr := host + ":" + port

	// init and start the server
	server := server.NewServer(addr, pgStorage, minioStorage)
	server.Run()
}
