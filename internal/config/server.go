package config

func (c *Config) IsAllowedToCreateNewChannels() bool {
	if c.yamlConfig == nil {
		panic(errYamlConfigNotDefined)
	}
	return c.yamlConfig.Server.AllowCreateNewChannels
}

func (c *Config) IsAllowedToPushServerMessages() bool {
	if c.yamlConfig == nil {
		panic(errYamlConfigNotDefined)
	}
	return c.yamlConfig.Server.AllowPushServerMessages
}
