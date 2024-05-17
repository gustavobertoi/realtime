package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/dtos"
)

func (h *Handler) CreateNewChannelHandler(c *gin.Context) {
	logger := h.Logger

	var dto = &dtos.CreateChannelDTO{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		logger.Errorf("error validating dto to create a channel, err: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "invalid body schema",
		})
		return
	}

	ch, err := channels.NewChannel(dto, h.Conf.PubSub.Consumer, h.Conf.PubSub.Producer)
	if err != nil {
		logger.Errorf("error creating channel, err: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	logger.Infof("channel %s created with success", ch.ID)

	c.IndentedJSON(http.StatusOK, ch)
}
