package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/channels"
	"github.com/open-source-cloud/realtime/config"
)

func SendServerMessageHandler(c *gin.Context) {
	logger := config.NewLogger("[POST] /channels/:channelId/messages")

	var dto = &channels.SendServerMessageDTO{}
	err := c.BindJSON(&dto)
	if err != nil {
		// TODO: improve validations
		logger.Errorf("error validating server message, err %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "invalid body schema",
		})
		return
	}

	channelID := c.Param("channelId")
	channel, err := conf.GetChannelByID(channelID)
	if err != nil {
		logger.Errorf("error getting channel by id %s, err: %v", channelID, err)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "channel not found",
		})
		return
	}

	clientId := "server-message"
	msg := channels.NewMessage(channelID, clientId, dto.Payload)

	if err := channel.BroadcastMessage(msg); err != nil {
		logger.Errorf("error broadcasting server message to all clients, err: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Error broadcasting message to all clients",
		})
		return
	}

	logger.Infof("message %s has been sent to all clients of channel %s", msg.ID, msg.ChannelID)

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("message %s has been sent to all clients", msg.ID),
	})
}
