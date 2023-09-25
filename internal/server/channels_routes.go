package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/internal/channels"
	"github.com/open-source-cloud/realtime/internal/config"
	"github.com/open-source-cloud/realtime/internal/server/drivers"
	"github.com/open-source-cloud/realtime/pkg/log"
	"github.com/sirupsen/logrus"
)

var systemLog = log.GetStaticInstance()

func channelById(c *gin.Context, conf *config.Config) {
	channelID := c.Param("channelId")
	channel, err := conf.GetChannelByID(channelID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Channel not found",
		})
		return
	}

	if channel.IsMaxOfConnections() {
		err := fmt.Errorf("maximum connection limit for this channel %s has been established", channelID)
		systemLog.Error(err.Error())
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	userAgent, ip := GetIPAndUserAgent(c.Request)
	client := channels.NewClient(&channels.CreateClientDTO{
		ID:        c.Query("clientId"),
		ChannelID: channel.ID,
		UserAgent: userAgent,
		IPAddress: ip,
	})
	channel.AddClient(client)

	upgradeConnection := c.Query("upgrade")
	if upgradeConnection == "" {
		upgradeConnection = "1"
	}

	logger := log.CreateWithContext("channels_routes.go", logrus.Fields{
		"channel_id":        channelID,
		"client_id":         client.ID,
		"client_user_agent": userAgent,
		"client_ip_address": ip,
	})

	if upgradeConnection == "0" {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "OK",
		})
		return
	}

	logger.Infof("client %s from channel %s has been created, upgrading connection to %s", client.ID, channelID, channel.Type)

	if channel.IsWebSocket() {
		drivers.WebSocket(c.Request, c.Writer, client, channel)
		return
	}

	if channel.IsSSE() {
		drivers.NewSSE(c, channel, client)
		return
	}
}

func pushServerMessage(c *gin.Context, conf *config.Config) {
	var dto *channels.SendServerMessageDTO
	err := c.BindJSON(&dto)
	if err != nil {
		// TODO: improve validations
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "invalid body schema",
		})
		return
	}

	channelID := c.Param("channelId")
	channel, err := conf.GetChannelByID(channelID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "channel not found",
		})
		return
	}

	clientId := "server-message"
	msg := channels.NewMessage(channelID, clientId, dto.Payload)

	if err := channel.BroadcastMessage(msg); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Error broadcasting message to all clients",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("message %s has been sent to all clients", msg.ID),
	})
}
