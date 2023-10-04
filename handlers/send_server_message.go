package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/channels"
)

func SendServerMessage(c *gin.Context) {
	var dto = &channels.SendServerMessageDTO{}
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
