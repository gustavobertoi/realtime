package config

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
