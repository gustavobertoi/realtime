package main

import (
	"github.com/open-source-cloud/realtime/internal/config"
	"github.com/open-source-cloud/realtime/internal/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

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
