# Client plugins

Client plugins must implement the `Handler` interface defined in the `plugin` package. The plugin should be registered
to the `PluginManager` in the `init` function by calling the `GetPluginManager().Register(commandName, plugin)`
or `Register(commandName, plugin)`.

## Plugin interface

```go
package example

type Handler interface {
	// Execute method is called when a client receives a command from the server. As an argument, it gets the data from the Payload sent from the server. 
	// The execute method returns the data, which is inserted into the Payload and sent to the server as a reply to the command.
	Execute(args interface{}) (interface{}, error)

	// GetMetadata returns the metadata of the client plugin.
	GetMetadata() plugin.Metadata
}
```

If you want or need to send a message in a plugin without prior request, there is a `connection.SendToServer(payload)`
method that allows you to do so. The only argument is a `Payload` object.

## Plugin example

```go
package example

import "github.com/detecc/deteccted/plugin"
import "log"

func init() {
	examplePlugin := &ExamplePlugin{}
	plugin.Register(examplePlugin.GetCmdName(), examplePlugin)
}

type ExamplePlugin struct {
	plugin.Handler
}

func (e ExamplePlugin) GetCmdName() string {
	return "/exampleCmd"
}

func (e ExamplePlugin) Execute(args interface{}) (interface{}, error) {
	log.Println(args)
	return "ping", nil
}

func (e ExamplePlugin) GetMetadata() plugin.Metadata {
	return plugin.Metadata{Type: plugin.PluginTypeClientServer}
}
```

## Compiling the plugin

The plugin is compiled using the following command:

```bash
go build -buildmode=plugin . 
```

It produces a file with `.so` format. The file's name (without the format) is then specified in the configuration file
under `plugins`. 