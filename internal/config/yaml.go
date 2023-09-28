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

func (c *Config) LoadConfigFromYaml() {
	filePath := getConfigFilePath()
	file, err := os.ReadFile(filePath)
	// NOTE: Means that file does not exists in container memory (skipping it)
	if err != nil {
		return
	}
	var schema *RootConfigDTO
	err = yaml.Unmarshal(file, &schema)
	// NOTE: Error reading config file and parsing it (throw error?)
	if err != nil {
		panic(err)
	}
	// TODO: Add a validation method to validate all nested props/structs
	c.rootConfig = schema
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
