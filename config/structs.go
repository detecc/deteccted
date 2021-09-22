package config

type Configuration struct {
	ServiceNodeIdentifier string
	Client                Client
}

type Client struct {
	Host         string   `fig:"host" validate:"required"`
	Port         int      `fig:"port" validate:"required"`
	AuthPassword string   `fig:"authPassword" validate:"required"`
	PingInterval int      `fig:"pingInterval" validate:"required"`
	PluginDir    string   `fig:"pluginDir"`
	Plugins      []string `fig:"plugins"`
}
