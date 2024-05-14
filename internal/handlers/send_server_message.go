package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/dtos"
)

func (h *Handler) SendServerMessageHandler(c *gin.Context) {
	logger := h.Logger

	var dto = &dtos.SendServerMessageDTO{}
	err := c.BindJSON(&dto)
	if err != nil {
		logger.Errorf("error validating server message, err %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "invalid body schema",
		})
		return
	}

	channelID := c.Param("channelId")
	channel := h.Conf.GetChannel(channelID)
	if channel == nil {
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
