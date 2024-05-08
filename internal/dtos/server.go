package dtos

type ServerConfig struct {
	AllowCreateNewChannels  bool `yaml:"allow_create_new_channels"`
	AllowPushServerMessages bool `yaml:"allow_push_server_messages"`
	AllowAllOrigins         bool `yaml:"allow_all_origins"`
	RenderChatHTML          bool `yaml:"render_chat_html"`
	RenderNotificationsHTML bool `yaml:"render_notifications_html"`
}
