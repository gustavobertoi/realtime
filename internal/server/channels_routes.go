package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/open-source-cloud/realtime/internal/channels"
	"github.com/open-source-cloud/realtime/internal/config"
	"github.com/open-source-cloud/realtime/pkg/log"
	"github.com/open-source-cloud/realtime/pkg/uuid"
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

	if clientStore.Count() >= channel.Config.MaxOfChannelConnections {
		err := fmt.Errorf("maximum connection limit for this channel %s has been established", channelID)
		systemLog.Error(err.Error())
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
	}

	clientID := c.Query("clientId")
	if clientID == "" {
		clientID = uuid.NewUUID()
	}
	userAgent, ip := GetIPAndUserAgent(c.Request)
	client := channels.NewClient(clientID, userAgent, ip, channelID)
	client.SetProducerAdapter(conf.GetClientProducerAdapter())
	if clientStore.Put(client); err != nil {
		systemLog.Errorf("error setting client %s from channel %s into store, details: %v", clientID, channelID, err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	logger := log.CreateWithContext("channels_routes.go", logrus.Fields{
		"channel_id": channelID,
		"client_id":  clientID,
	})

	logger.Print("client created upgrading connection to websocket")

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
		clientStore.Delete(client)
	}()

	// Write messages
	go func() {
		ch := client.GetChan()
		if channel.Subscribe(client); err != nil {
			logger.Panicf("error subscribing client on channel")
			return
		}
		for {
			msg := <-ch
			logger.Infof("writing message %s to buffer", msg.ID)
			err := writeMessageToBuffer(msg, client, conn)
			if err != nil {
				logger.Errorf("error writing message %s to buffer, details: %v", msg.ID, err)
			} else {
				logger.Infof("message %s has been sent to client", msg.ID)
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
		logger.Printf("broadcasting message %s to all clients", msg.ID)
		if messageType == websocket.TextMessage {
			for _, err := range channel.BroadcastToAllClients(msg) {
				logger.Errorf("error broadcasting message %s, details: %v", msg.ID, err)
			}
		}
	}
}

func writeMessageToBuffer(msg *channels.Message, client *channels.Client, conn *websocket.Conn) error {
	// TODO: Write self messages?
	if msg.ClientID != client.ID {
		msgStr, err := msg.ToJSON()
		if err != nil {
			return err
		}
		return conn.WriteMessage(websocket.TextMessage, []byte(msgStr))
	}
	return nil
}
