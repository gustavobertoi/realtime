package drivers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/open-source-cloud/realtime/internal/channels"
)

func WebSocket(request *http.Request, writer gin.ResponseWriter, client *channels.Client, channel *channels.Channel) {
	if err := channel.Subscribe(client); err != nil {
		panic(err)
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		panic(err)
	}

	msgChan := client.MessageChan()
	go func() {
		for {
			msg := <-msgChan
			if msg.ClientID != client.ID {
				msgStr, err := msg.MessageToJSON()
				if err != nil {
					panic(err)
				}
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msgStr)); err != nil {
					panic(err)
				}
			}
		}
	}()

	// Closes WS client connection
	defer func() {
		conn.Close()
		channel.DeleteClient(client)
	}()

	// Read WS messages
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}
		msg := channels.NewMessage(channel.ID, client.ID, string(payload))
		if messageType == websocket.TextMessage {
			if err := channel.BroadcastMessage(msg); err != nil {
				break
			}
		}
	}
}
