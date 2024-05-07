package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gustavobertoi/realtime/channels"
	"github.com/gustavobertoi/realtime/config"
)

func UpgradeChannelConnectionHandler(c *gin.Context) {
	logger := config.NewLogger("[POST] /api/v1/channels")

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
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	userAgent := c.Request.UserAgent()

	client := channels.NewClient(&channels.CreateClientDTO{
		ID:        c.Query("clientId"),
		ChannelID: channel.ID,
		UserAgent: userAgent,
		IPAddress: ip,
	})

	if channel.HasClient(client) {
		logger.Warnf("client %s already connected with channel %s returning 409", client.ID, channelID)
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": fmt.Sprintf("client %s is already connected on channel, please disconnect it before connect it again", client.ID),
		})
		return
	}

	channel.AddClient(client)

	upgradeConnection := c.Query("upgrade")
	if upgradeConnection == "" {
		upgradeConnection = "1"
	}

	if upgradeConnection == "0" {
		channel.DeleteClient(client)
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "OK",
		})
		return
	}

	logger.Infof("client %s from channel %s has been created, upgrading connection to %s", client.ID, channelID, channel.Type)

	if err := channel.Subscribe(client); err != nil {
		logger.Errorf("error subscribing client on adapter: %v", err)
		panic(err)
	}

	serverConf := conf.GetServerConfig()

	switch channel.Type {
	case channels.WebSocket:
		WebSocketHandler(c, serverConf, channel, client, logger)
	case channels.ServerSentEvents:
		ServerSentEventsHandler(c, serverConf, channel, client, logger)
	}

	channel.DeleteClient(client)

	logger.Warnf("the channel %s connection with the client %s has ended", channel.ID, client.ID)
}
