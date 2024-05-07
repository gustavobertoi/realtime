package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustavobertoi/realtime/channels"
	"github.com/gustavobertoi/realtime/config"
)

func CreateNewChannelHandler(c *gin.Context) {
	logger := config.NewLogger("[POST] /channels")

	svConf := conf.GetServerConfig()

	if !svConf.AllowCreateNewChannels {
		logger.Errorf("not allowed to create new channels")
		c.IndentedJSON(http.StatusForbidden, gin.H{
			"message": "server are not allowed to create new channels",
		})
		return
	}

	var dto = &channels.CreateChannelDTO{}
	err := c.BindJSON(&dto)
	if err != nil {
		// TODO: improve validations
		logger.Errorf("error validating dto to create a channel, err: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "invalid body schema",
		})
		return
	}

	ch, err := conf.CreateChannel(dto, nil)
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
