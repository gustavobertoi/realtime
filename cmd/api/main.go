package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gustavobertoi/realtime/internal/config"
	"github.com/gustavobertoi/realtime/internal/handlers"
)

func main() {
	c, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	if c.Server.AllowAllOrigins {
		r.Use(cors.Default())
	}

	apiV1 := r.Group("/api/v1")

	apiV1.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "UP"})
	})

	handler := handlers.NewHandler(c)

	// Channels
	apiV1.POST("/channels", handler.CreateNewChannelHandler)
	apiV1.GET("/channels/:channelId", handler.UpgradeChannelConnectionHandler)

	apiV1.POST("/channels/:channelId/messages", handler.UpgradeChannelConnectionHandler)

	if err := r.Run(c.Port()); err != nil {
		log.Fatal(err.Error())
	}
}
