package config

import (
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type (
	YamlServerDTO struct {
		AllowCreateNewChannels  bool `yaml:"allow_create_new_channels"`
		AllowPushServerMessages bool `yaml:"allow_push_server_messages"`
	}
	YamlPubSubRedisDTO struct {
		URL string `yaml:"url"`
	}
	YamlPubSubDTO struct {
		Driver string              `yaml:"driver"`
		Redis  *YamlPubSubRedisDTO `yaml:"redis"`
	}
	YamlChannelConfigDTO struct {
		MaxOfConnections int `yaml:"max_of_connections"`
	}
	YamlChannelDTO struct {
		ID     string                `yaml:"id"`
		Type   string                `yaml:"type"`
		Name   string                `yaml:"name"`
		Config *YamlChannelConfigDTO `yaml:"config"`
	}
	YamlConfigRootDTO struct {
		Server   *YamlServerDTO             `yaml:"server"`
		PubSub   *YamlPubSubDTO             `yaml:"pubsub"`
		Channels map[string]*YamlChannelDTO `yaml:"channels"`
	}
)

func (c *Config) LoadConfigYaml() error {
	filePath, err := loadConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var schema *YamlConfigRootDTO
	err = yaml.Unmarshal(file, &schema)
	if err != nil {
		return err
	}

	c.yamlConfig = schema

	if err := c.createChannelsFromConfig(); err != nil {
		return err
	}

	return nil
}

func loadConfigFilePath() (string, error) {
	envFilePath := os.Getenv("CONFIG_FOLDER_PATH")
	if envFilePath != "" {
		return envFilePath, nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fileName := "config.yaml"
	resourcesFolder := "realtime"
	filePath := path.Join(pwd, resourcesFolder, fileName)

	return filePath, nil
}
