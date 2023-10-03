package config

import (
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type (
	ServerDTO struct {
		AllowCreateNewChannels  bool `yaml:"allow_create_new_channels"`
		AllowPushServerMessages bool `yaml:"allow_push_server_messages"`
		RenderChatHTML          bool `yaml:"render_chat_html"`
		RenderNotificationsHTML bool `yaml:"render_notifications_html"`
	}
	PubSubRedisDTO struct {
		URL string `yaml:"url"`
	}
	PubSubDTO struct {
		Driver string          `yaml:"driver"`
		Redis  *PubSubRedisDTO `yaml:"redis"`
	}
	ChannelConfigDTO struct {
		MaxOfConnections int `yaml:"max_of_connections"`
	}
	ChannelDTO struct {
		ID     string            `yaml:"id"`
		Type   string            `yaml:"type"`
		Name   string            `yaml:"name"`
		Config *ChannelConfigDTO `yaml:"config"`
	}
	RootConfigDTO struct {
		Server   *ServerDTO             `yaml:"server"`
		PubSub   *PubSubDTO             `yaml:"pubsub"`
		Channels map[string]*ChannelDTO `yaml:"channels"`
	}
)

func getConfigFilePath() string {
	envFilePath := os.Getenv("CONFIG_FOLDER_PATH")
	if envFilePath != "" {
		return envFilePath
	}
	pwd, _ := os.Getwd()
	fileName := "config.yaml"
	resourcesFolder := "realtime"
	filePath := path.Join(pwd, resourcesFolder, fileName)
	return filePath
}

func (c *Config) LoadConfigFromYaml() error {
	filePath := getConfigFilePath()
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var schema *RootConfigDTO
	err = yaml.Unmarshal(file, &schema)
	if err != nil {
		return err
	}
	c.rootConfig = schema
	return nil
}

func (c *Config) GetServerConfig() *ServerDTO {
	return c.rootConfig.Server
}

func (c *Config) GetPubSubConfig() *PubSubDTO {
	return c.rootConfig.PubSub
}

func (c *Config) GetChannelsConfig() map[string]*ChannelDTO {
	return c.rootConfig.Channels
}
