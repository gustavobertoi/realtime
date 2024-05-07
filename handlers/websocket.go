package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gustavobertoi/realtime/channels"
	"github.com/gustavobertoi/realtime/config"
)

func WebSocketHandler(c *gin.Context, serverConf *config.Server, channel *channels.Channel, client *channels.Client, logger *config.Logger) {
	// TODO: Improve websocket upgrade connections (read those infos from server config)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return serverConf.AllowAllOrigins },
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
