# Client plugins

Client plugins must implement the `Handler` interface defined in the `plugin` package. The plugin should be registered
to the `PluginManager` in the `init` function by calling the `GetPluginManager().Register(commandName, plugin)`
or `Register(commandName, plugin)`.

## Plugin interface

```golang
type Handler interface {
// Execute method is called when a client receives a command from the server. As an argument, it gets the data from the Payload sent from the server. 
// The execute method returns the data, which is inserted into the Payload and sent to the server as a reply to the command.
Execute(args interface{}) (interface{}, error)
}
```

## Plugin example

## Compiling the plugin

The plugin is compiled using the following command:

```bash
GO111MODULE=off go build -buildmode=plugin . 
```

It produces a file with `.so` format. The file's name (without the format) is then specified in the configuration file
under `plugins`. 