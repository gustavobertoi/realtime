package dtos

type (
	RedisConfig struct {
		URL string `yaml:"url"`
	}
	UpstashKafkaConfig struct {
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		ServerAddr string `yaml:"server_addr"`
		Topic      string `yaml:"topic"`
		GroupId    string `yaml:"group_id"`
	}
	PubSubConfig struct {
		Driver       string              `yaml:"driver"`
		Redis        *RedisConfig        `yaml:"redis"`
		UpstashKafka *UpstashKafkaConfig `yaml:"upstash_kafka"`
	}
)
