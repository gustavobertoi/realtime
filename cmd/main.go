package main

import (
	"github.com/open-source-cloud/realtime/internal/config"
	"github.com/open-source-cloud/realtime/internal/server"
	"github.com/open-source-cloud/realtime/pkg/log"
)

func main() {
	logger := log.GetStaticInstance()

	logger.Print("Creating config")

	config := config.NewConfig()

	logger.Print("Creating server")
	server := server.NewServer(config)

	logger.Print("Starting http and ws server")
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
