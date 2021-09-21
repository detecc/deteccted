package plugin

import (
	"fmt"
	cache2 "github.com/patrickmn/go-cache"
	"log"
	"plugin"
	"github.com/detecc/deteccted/cache"
	"github.com/detecc/deteccted/config"
)

type Handler interface {
	// Execute method is called when a client receives a command from the server. As an argument, it gets the data attribute from the Payload.
	//The execute method returns the Payload ready to be sent as a reply to the server's command request.
	Execute(args interface{}) (interface{}, error)
}

// Register the plugin in the client.
func Register(name string, action Handler) {
	log.Println("Adding plugin to cache", name, action)
	cache.Memory().Set(name, action, cache2.NoExpiration)
}


// LoadPlugins load all the plugins specified in the plugin list. Each plugin should have a unique command.
func LoadPlugins() {
	log.Println("Loading plugins..")
	server := config.GetClientConfiguration().Client

	for _, pluginFromList := range server.Plugins {
		log.Println("Loading: ", pluginFromList)
		_, err := plugin.Open(fmt.Sprintf("%s/%s.so", server.PluginDir, pluginFromList))
		if err != nil {
			fmt.Println("error loading plugin", err)
			continue
		}
	}
}
