package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/config"
	"github.com/open-source-cloud/realtime/handlers"
)

func main() {
	c := config.GetConfig()
	if err := c.LoadConfigFromYaml(); err != nil {
		log.Fatal(err.Error())
	}
	if err := c.CreateChannelsFromConfig(); err != nil {
		log.Fatal(err.Error())
	}

	r := gin.New()
	r.Use(gin.Recovery())

	svConf := c.GetServerConfig()

	if svConf.AllowAllOrigins {
		r.Use(cors.Default())
	}

	if svConf.RenderChatHTML {
		r.Static("/web/chat", "./web/chat")
	}

	if svConf.RenderNotificationsHTML {
		r.Static("/web/notifications", "./web/notifications")
	}

	apiV1 := r.Group("/api/v1")

	// Channels
	apiV1.POST("/channels", handlers.CreateNewChannelHandler)
	apiV1.GET("/channels/:channelId", handlers.UpgradeChannelConnectionHandler)

	apiV1.POST("/channels/:channelId/messages", handlers.SendServerMessageHandler)

	if err := r.Run(c.GetPort()); err != nil {
		log.Fatal(err.Error())
	}
}
