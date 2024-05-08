package config

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/dtos"
	"github.com/gustavobertoi/realtime/pkg/pubsub"
	"github.com/gustavobertoi/realtime/pkg/store"
	"gopkg.in/yaml.v2"
)

var (
	ErrYamlConfigNotDeclared = errors.New("yaml config file not declared as env var")
)

type (
	YamlConfig struct {
		Server   *dtos.ServerConfig
		PubSub   *dtos.PubSubConfig
		Channels map[string]*dtos.CreateChannelDTO
	}
	Config struct {
		Debug    bool
		Store    *store.MemoryStore
		Server   *dtos.ServerConfig
		PubSub   *pubsub.PubSub
		Channels map[string]*channels.Channel
	}
)

var conf *Config

func GetConfig() (*Config, error) {
	if conf != nil {
		return conf, nil
	}

	memoryPubSub, err := pubsub.NewMemoryAdapter()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Debug: os.Getenv("APP_DEBUG") == "1",
		Store: store.NewMemoryStore(),
		Server: &dtos.ServerConfig{
			AllowCreateNewChannels:  true,
			AllowPushServerMessages: true,
			RenderChatHTML:          false,
			RenderNotificationsHTML: false,
		},
		PubSub:   memoryPubSub,
		Channels: make(map[string]*channels.Channel),
	}

	yamlConfig, err := getConfigFromYaml()
	if err != nil {
		log.Printf("error reading yaml config (using default one): %s", err)
		return config, nil
	}

	config.Server = yamlConfig.Server

	pubsub, err := pubsub.NewPubSub(context.Background(), yamlConfig.PubSub)
	if err != nil {
		return nil, err
	}

	for _, channelDTO := range yamlConfig.Channels {
		channel, err := channels.NewChannel(channelDTO, pubsub.Consumer, pubsub.Producer)
		if err != nil {
			return nil, err
		}
		config.SetChannel(channel)
	}

	return config, nil
}

func (c *Config) Port() string {
	return ":4567"
}

func (c *Config) GetChannel(id string) *channels.Channel {
	channel, ok := c.Channels[id]
	if !ok {
		return nil
	}
	return channel
}

func (c *Config) SetChannel(channel *channels.Channel) {
	c.Channels[channel.ID] = channel
}

func getConfigFromYaml() (*YamlConfig, error) {
	filePath := os.Getenv("CONFIG_FOLDER_PATH")
	if filePath != "" {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		var schema YamlConfig
		err = yaml.Unmarshal(file, &schema)
		if err != nil {
			return nil, err
		}
		return &schema, nil
	}
	return nil, ErrYamlConfigNotDeclared
}
