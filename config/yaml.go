package config

import (
	"os"

	"github.com/open-source-cloud/realtime/channels"
	"gopkg.in/yaml.v2"
)

type (
	Server struct {
		AllowCreateNewChannels  bool `yaml:"allow_create_new_channels"`
		AllowPushServerMessages bool `yaml:"allow_push_server_messages"`
		AllowAllOrigins         bool `yaml:"allow_all_origins"`
		RenderChatHTML          bool `yaml:"render_chat_html"`
		RenderNotificationsHTML bool `yaml:"render_notifications_html"`
	}
	PubSubRedis struct {
		URL string `yaml:"url"`
	}
	PubSubKafka struct {
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		ServerAddr string `yaml:"server_addr"`
		Topic      string `yaml:"topic"`
		GroupId    string `yaml:"group_id"`
	}
	PubSub struct {
		Driver string       `yaml:"driver"`
		Redis  *PubSubRedis `yaml:"redis"`
		Kafka  *PubSubKafka `yaml:"kafka"`
	}
	RootConfigDTO struct {
		Server   *Server                               `yaml:"server"`
		PubSub   *PubSub                               `yaml:"pubsub"`
		Channels map[string]*channels.CreateChannelDTO `yaml:"channels"`
	}
)

func (c *Config) LoadConfigFromYaml() error {
	filePath := os.Getenv("CONFIG_FOLDER_PATH")
	if filePath != "" {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		var schema = &RootConfigDTO{}
		err = yaml.Unmarshal(file, &schema)
		if err != nil {
			return err
		}
		// TODO: Add a validation schema
		c.rootConfig = schema
	}
	return nil
}

func (c *Config) GetServerConfig() *Server {
	return c.rootConfig.Server
}

func (c *Config) GetPubSubConfig() *PubSub {
	return c.rootConfig.PubSub
}

func (c *Config) GetChannelsConfig() map[string]*channels.CreateChannelDTO {
	return c.rootConfig.Channels
}
