package plugin

const (
	PluginTypeClientServer = "clientServer"
	PluginTypeClientOnly   = "clientOnly"
)

type (
	// Handler is the interface for the Plugin.
	Handler interface {
		// Execute method is called when a client receives a command from the server.
		// The arguments of the method is Payload's data attribute, incoming from the server. The response should be data, ready to
		// be sent back to the server for processing and an error, if one occurred.
		Execute(args interface{}) (interface{}, error)

		// GetMetadata returns the metadata of the client plugin.
		GetMetadata() Metadata
	}

	// Metadata is the metadata object for the plugin and determines the behaviour of the plugin interaction with the client.
	Metadata struct {
		Type string
	}
)
