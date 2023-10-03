package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/open-source-cloud/realtime/channels"
	"github.com/open-source-cloud/realtime/config"
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
	channel.AddClient(client)

	upgradeConnection := c.Query("upgrade")
	if upgradeConnection == "" {
		upgradeConnection = "1"
	}

	if upgradeConnection == "0" {
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

	switch channel.Type {
	case channels.WebSocket:
		ws(c, channel, client, logger)
	case channels.ServerSentEvents:
		sse(c, channel, client, logger)
	}

	channel.DeleteClient(client)

	logger.Warnf("the channel %s connection with the client %s has ended", channel.ID, client.ID)
}

func ws(c *gin.Context, channel *channels.Channel, client *channels.Client, logger *config.Logger) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("error upgrading connection to websocket, err: %v", err)
		panic(err)
	}

	msgChan := client.MessageChan()
	go func() {
		for {
			msg := <-msgChan
			if msg.ClientID != client.ID {
				logger.Infof("serializing and writing msg %s to buffer", msg.ID)
				msgStr, err := msg.MessageToJSON()
				if err != nil {
					logger.Errorf("error serializing msg %s to json, err: %v", msg.ID, err)
					break
				}
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msgStr)); err != nil {
					logger.Errorf("error writing msg %s on buffer, err: %v", msg.ID, err)
					break
				}
				logger.Infof("msg %s was written to buffer for client %s", msg.ID, client.ID)
			} else {
				logger.Warnf("not writing self msg %s to this client %s", msg.ID, client.ID)
			}
		}
	}()

	// Closes WS client connection
	defer conn.Close()

	// Read WS messages
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			logger.Errorf("error reading message from buffer, err: %v", err)
			break
		}
		msg := channels.NewMessage(channel.ID, client.ID, string(payload))
		logger.Infof("sending %s msg to all clients", msg.ID)
		if messageType == websocket.TextMessage {
			if err := channel.BroadcastMessage(msg); err != nil {
				logger.Errorf("error broadcasting msg %s to clients, err: %v", msg.ID, err)
				break
			}
			logger.Infof("msg %s has been sent to all client", msg.ID)
		}
	}
}

func sse(c *gin.Context, channel *channels.Channel, client *channels.Client, logger *config.Logger) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	msgChan := client.MessageChan()
	clientGone := c.Writer.CloseNotify()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-clientGone:
				return
			case msg := <-msgChan:
				logger.Infof("serializing and writing msg %s to client %s", msg.ID, client.ID)
				if msg.ClientID != client.ID {
					msgStr, err := msg.MessageToJSON()
					if err != nil {
						logger.Errorf("error serializing msg %s to json, err: %v", msg.ID, err)
						break
					}
					message := fmt.Sprintf("data: %s\n\n", msgStr)
					_, err = c.Writer.WriteString(message)
					if err != nil {
						logger.Errorf("error writing msg %s on buffer, err: %v", msg.ID, err)
						break
					}
					c.Writer.Flush()
					logger.Infof("msg %s was written to buffer for client %s", msg.ID, client.ID)
					time.Sleep(1 * time.Second)
				} else {
					logger.Warnf("not writing self msg %s to client %s", msg.ID, client.ID)
				}
			}
		}
	}()

	wg.Wait()
}
