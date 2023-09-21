package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/open-source-cloud/realtime/internal/channels"
	"github.com/open-source-cloud/realtime/internal/config"
	"github.com/open-source-cloud/realtime/pkg/log"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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
	}

	userAgent, ip := GetIPAndUserAgent(c.Request)
	client := channels.NewClient(&channels.CreateClientDTO{
		ID:        c.Query("clientId"),
		ChannelID: channel.ID,
		IPAddress: ip,
		UserAgent: userAgent,
	})

	if channel.Store.Put(client); err != nil {
		systemLog.Errorf("error setting client %s from channel %s into store, details: %v", client.ID, channelID, err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	logger := log.CreateWithContext("channels_routes.go", logrus.Fields{
		"channel_id": channelID,
		"client_id":  client.ID,
	})

	logger.Infof("client %s from channel %s has been created, upgrading connection to websocket", client.ID, channelID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("failed to set websocket upgrade: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer func() {
		logger.Print("deleting client from channel")
		conn.Close()
		channel.Store.Delete(client)
	}()

	// Write messages
	go func() {
		if channel.Subscribe(client); err != nil {
			logger.Panicf("error subscribing client on channel")
			return
		}
		for {
			msg := <-client.MessageChan()
			// TODO: Write self msgs?
			if msg.ClientID != client.ID {
				logger.Infof("writing message %s to buffer", msg.ID)
				err := writeMessageToBuffer(msg, client, conn)
				if err != nil {
					logger.Errorf("error writing message %s to buffer, details: %v", msg.ID, err)
					continue
				}
				logger.Infof("message %s has been sent to client %s", msg.ID, client.ID)
			}
		}
	}()

	// Read messages
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("error reading message from buffer: %v", err)
			break
		}
		msg := channels.NewMessage(channelID, client.ID, string(payload))
		logger.Printf("channel %s is broadcasting message %s from client %s to all clients", channelID, msg.ID, client.ID)
		if messageType == websocket.TextMessage {
			err := channel.BroadcastMessage(msg)
			if err != nil {
				logger.Errorf("error broadcasting message %s, details: %v", msg.ID, err)
			}
		}
	}
}

func writeMessageToBuffer(msg *channels.Message, client *channels.Client, conn *websocket.Conn) error {
	msgStr, err := msg.MessageToJSON()
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(msgStr))
}
