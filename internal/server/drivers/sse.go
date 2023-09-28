package drivers

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-source-cloud/realtime/internal/channels"
	"github.com/sirupsen/logrus"
)

func NewSSE(c *gin.Context, channel *channels.Channel, client *channels.Client, logger *logrus.Entry) {
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
