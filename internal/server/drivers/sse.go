package drivers

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/internal/channels"
)

func NewSSE(c *gin.Context, channel *channels.Channel, client *channels.Client) {
	if err := channel.Subscribe(client); err != nil {
		panic(err)
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	msgChan := client.MessageChan()
	clientGone := c.Writer.CloseNotify()
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()
		for {
			select {
			case <-clientGone:
				return
			case msg := <-msgChan:
				msgStr, err := msg.MessageToJSON()
				if err != nil {
					break
				}
				message := fmt.Sprintf("data: %s\n\n", msgStr)
				c.Writer.WriteString(message)
				c.Writer.Flush()
				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg.Wait()
	channel.DeleteClient(client)
}
