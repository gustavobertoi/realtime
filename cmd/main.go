package main

import (
	"log"

	"github.com/open-source-cloud/realtime/internal/config"
	"github.com/open-source-cloud/realtime/internal/server"
)

func main() {
	config := config.NewConfig()
	if err := config.LoadConfigYaml(); err != nil {
		log.Fatal(err)
	}
	server := server.NewServer(config)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
