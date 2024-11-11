package main

import (
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/server"
)

func main() {
	// craft the server address using env vars
	host := config.Envs.Host
	port := config.Envs.Port
	addr := host + ":" + port

	// init and start the server
	server := server.NewServer(addr)
	server.Run()
}
