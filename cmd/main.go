package main

import (
	"log"

	"github.com/open-source-cloud/realtime/internal/config"
	"github.com/open-source-cloud/realtime/internal/server"
)

func main() {
	log.Print("Creating config")

	config := config.NewConfig()

	log.Print("Creating server")
	server := server.NewServer(config)

	log.Print("Starting http and ws server")
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
