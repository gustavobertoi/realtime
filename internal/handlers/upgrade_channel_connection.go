package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/dtos"
)

func (h *Handler) UpgradeChannelConnectionHandler(c *gin.Context) {
	channelID := c.Param("channelId")
	channel := h.Conf.GetChannel(channelID)
	if channel == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Channel not found",
		})
		return
	}

	if channel.IsMaxOfConnections() {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{
			"message": fmt.Sprintf("maximum connection limit for this channel %s has been established", channelID),
		})
		return
	}

	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	userAgent := c.Request.UserAgent()

	client := channels.NewClient(&dtos.CreateClientDTO{
		ID:        c.Query("clientId"),
		ChannelID: channel.ID,
		UserAgent: userAgent,
		IPAddress: ip,
	})

	if channel.HasClient(client) {
		h.Logger.Warnf("client %s already connected with channel %s returning 409", client.ID, channelID)
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

	h.Logger.Infof("client %s from channel %s has been created, upgrading connection to %s", client.ID, channelID, channel.Type)

	if err := channel.Subscribe(client); err != nil {
		h.Logger.Errorf("error subscribing client on adapter: %v", err)
		panic(err)
	}

	switch channel.Type {
	case channels.WebSocket:
		WebSocketHandler(c, h.Conf, channel, client, h.Logger)
	case channels.ServerSentEvents:
		ServerSentEventsHandler(c, h.Conf, channel, client, h.Logger)
	}

	channel.DeleteClient(client)

	h.Logger.Warnf("the channel %s connection with the client %s has ended", channel.ID, client.ID)
}
