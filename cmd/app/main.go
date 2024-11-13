package main

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/db"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/server"
)

func main() {
	// extracting dsn from env var
	dsn := config.Envs.DB_DSN
	if dsn == "" {
		panic("Database DSN not found")
	}

	// get new db storage
	pgStorage := db.NewPostgreSQLStorage(dsn)

	// craft the server address using env vars
	host := config.Envs.Host
	port := config.Envs.Port
	addr := host + ":" + port

	// init and start the server
	server := server.NewServer(addr, pgStorage)
	server.Run()
}
