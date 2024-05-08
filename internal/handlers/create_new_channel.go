package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/dtos"
	"github.com/gustavobertoi/realtime/pkg/logs"
)

func (h *Handler) CreateNewChannelHandler(c *gin.Context) {
	logger := logs.NewLogger("[POST] /channels")

	if !h.Conf.Server.AllowCreateNewChannels {
		logger.Errorf("not allowed to create new channels")
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "server are not allowed to create new channels",
		})
		return
	}

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
