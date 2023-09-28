package config

import (
	"fmt"
	"os"

	"github.com/open-source-cloud/realtime/pkg/store"
)

type Config struct {
	appDebug     bool
	port         int
	channelStore *store.MemoryStore
	rootConfig   *RootConfigDTO
}

func NewConfig() *Config {
	appDebug := os.Getenv("APP_DEBUG")
	c := &Config{
		appDebug:     appDebug == "1",
		port:         8080,
		channelStore: store.NewMemoryStore(),
		rootConfig: &RootConfigDTO{
			Server: &ServerDTO{
				AllowCreateNewChannels:  true,
				AllowPushServerMessages: true,
				RenderChatHTML:          false,
				RenderNotificationsHTML: false,
			},
			PubSub: &PubSubDTO{
				Driver: memoryDriver,
			},
			Channels: make(map[string]*ChannelDTO),
		},
	}
	return c
}

func (c *Config) GetPort() string {
	return fmt.Sprintf(":%d", c.port)
}

func (c *Config) IsAppOnDebugMode() bool {
	return c.appDebug
}
