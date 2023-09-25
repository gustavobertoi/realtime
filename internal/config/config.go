package config

import (
	"os"
	"path"

	"github.com/open-source-cloud/realtime/pkg/store"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int

	channelStore *store.MemoryStore
	yamlConfig   *YamlConfigRootDTO
}

func NewConfig() *Config {
	c := &Config{
		Port:         8080,
		channelStore: store.NewMemoryStore(),
		yamlConfig:   nil,
	}
	return c
}

func (c *Config) LoadConfigYaml() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := "config.yaml"
	resourcesFolder := "resources"
	filePath := path.Join(pwd, resourcesFolder, fileName)

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
