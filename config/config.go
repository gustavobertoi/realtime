package config

import (
	"fmt"
	"os"

	"github.com/open-source-cloud/realtime/channels"
	"github.com/open-source-cloud/realtime/pkg/store"
)

type Config struct {
	appDebug     bool
	port         int
	channelStore *store.MemoryStore
	rootConfig   *RootConfigDTO
}

var conf *Config

func GetConfig() *Config {
	if conf != nil {
		return conf
	}
	appDebug := os.Getenv("APP_DEBUG")
	conf = &Config{
		appDebug:     appDebug == "1",
		port:         8080,
		channelStore: store.NewMemoryStore(),
		rootConfig: &RootConfigDTO{
			Server: &Server{
				AllowCreateNewChannels:  true,
				AllowPushServerMessages: true,
				RenderChatHTML:          false,
				RenderNotificationsHTML: false,
			},
			PubSub: &PubSub{
				Driver: memoryDriver,
			},
			Channels: make(map[string]*channels.CreateChannelDTO),
		},
	}
	return conf
}

func (c *Config) GetPort() string {
	return fmt.Sprintf(":%d", c.port)
}

func (c *Config) IsAppOnDebugMode() bool {
	return c.appDebug
}
