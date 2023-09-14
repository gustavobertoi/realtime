package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/open-source-cloud/realtime/internal/channels"
	"github.com/open-source-cloud/realtime/internal/config"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func channelById(c *gin.Context, conf *config.Config) {
	channelID := c.Param("channel_id")
	if channelID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "missing channel_id",
		})
		return
	}

	channel, err := conf.GetChannelByID(channelID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "channel not found",
		})
		return
	}

	client, err := channel.CreateClient(GetIPAndUserAgent(c.Request))
	if err != nil {
		log.Printf("error adding client %s to channel %s, err: %v", client.ID, channel.ID, err)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	channel.SetAdapter(conf.GetChannelAdapter())
	client.SetAdapter(conf.GetClientAdapter())

	log.Printf("client %s has been connected to channel %s", client.ID, channel.ID)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("failed to set websocket upgrade: ", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer func() {
		log.Printf("Deleting client %s from channel %s", client.ID, channelID)
		conn.Close()
		channel.DeleteClient(client.ID)
	}()

	// Write messages
	go func() {
		ch := make(chan *channels.Message)
		err := channel.Subscribe(client)
		if err != nil {
			log.Printf("error subscribing to channel %s from client %s", channelID, client.ID)
		}
		for {
			msg := <-ch
			msgStr, err := msg.ToJSON()
			if err != nil {
				log.Printf("error deserializing msg %s from client %s - channel %s to json", msg.ID, msg.ClientID, msg.ChannelID)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, []byte(msgStr))
			if err != nil {
				log.Printf("error writing msg %s to client %s from channel %s", msg.ID, client.ID, channelID)
				return
			}
			log.Printf("channel %s sended a msg %s to client %s", channelID, msg.ID, client.ID)
			err = client.DeleteMessage(msg)
			if err != nil {
				log.Printf("error deleting msg %s from store of channel %s client %s", msg.ID, msg.ChannelID, msg.ClientID)
				return
			}
		}
	}()

	// Read messages
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading client %s msg from channel %s", client.ID, channelID)
			return
		}
		msg := channels.NewMessage(channelID, client.ID, payload)
		log.Printf("client %s is broadcasting a msg %s to channel %s server", client.ID, msg.ID, channelID)
		if messageType == websocket.TextMessage {
			channel.BroadcastAll(msg)
		}
	}
}
