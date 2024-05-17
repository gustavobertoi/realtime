package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/config"
	"github.com/gustavobertoi/realtime/pkg/logs"
)

func WebSocketHandler(c *gin.Context, conf *config.Config, channel *channels.Channel, client *channels.Client, logger *logs.Logger) {
	// TODO: Improve websocket upgrade connections (read those infos from server config)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return conf.Server.AllowAllOrigins
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("error upgrading connection to websocket, err: %v", err)
		panic(err)
	}

	msgChan := client.MessageChan()
	errChan := make(chan error)

	defer close(msgChan)
	defer close(errChan)

	go func() {
		for {
			msg := <-msgChan
			if err := writeWsMessage(conn, client, msg, logger); err != nil {
				errChan <- err
				return
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if err := readWsMessage(conn, channel, client, logger); err != nil {
				errChan <- err
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()

	for {
		if err := <-errChan; err != nil {
			logger.Error(err)
			if err := conn.Close(); err != nil {
				logger.Errorf("error closing connection, err: %v", err)
				return
			}
		}
	}
}

func readWsMessage(conn *websocket.Conn, channel *channels.Channel, client *channels.Client, logger *logs.Logger) error {
	messageType, payload, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("error reading message from buffer, err: %v", err)
	}
	msg := channels.NewMessage(channel.ID, client.ID, string(payload))
	logger.Infof("sending %s msg to all clients", msg.ID)
	if messageType == websocket.TextMessage {
		if err := channel.BroadcastMessage(msg); err != nil {
			return fmt.Errorf("error broadcasting msg %s to clients, err: %v", msg.ID, err)
		}
		logger.Infof("msg %s has been sent to all client", msg.ID)
	}
	return nil
}

func writeWsMessage(conn *websocket.Conn, client *channels.Client, msg *channels.Message, logger *logs.Logger) error {
	if msg.ClientID == client.ID {
		logger.Warnf("msg %s was sent by the same client %s, skipping", msg.ID, client.ID)
		return nil
	}
	logger.Infof("serializing and writing msg %s to buffer", msg.ID)
	msgStr, err := msg.MessageToJSON()
	if err != nil {
		return fmt.Errorf("error serializing msg %s to json, err: %v", msg.ID, err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msgStr)); err != nil {
		return fmt.Errorf("error writing msg %s on buffer, err: %v", msg.ID, err)
	}
	logger.Infof("msg %s was written to buffer for client %s", msg.ID, client.ID)
	return nil
}
