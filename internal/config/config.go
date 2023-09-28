package config

import (
	"fmt"

	"github.com/open-source-cloud/realtime/pkg/store"
)

type Config struct {
	port         int
	channelStore *store.MemoryStore
	yamlConfig   *YamlConfigRootDTO
}

func NewConfig() *Config {
	c := &Config{
		port:         8080,
		channelStore: store.NewMemoryStore(),
		yamlConfig:   nil,
	}
	return c
}

func (c *Config) GetPort() string {
	return fmt.Sprintf(":%d", c.port)
}
