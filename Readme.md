# 🎯 Deteccted 

Deteccted is a customizable, 🔌 plugin-based client for the ⚡[**Detecctor**](https://github.com/detecc/detecctor).

## ⚙ Configuration

An example configuration file for the client:

```yaml
serviceNodeIdentifier: "yourServiceNodeId"
client:
  host: localhost
  port: 7777
  authPassword: yourPassword
  pluginDir: "/usr/deteccted/plugins"
  plugins:
    - "examplePlugin"
```

## 🔌 Plugins

Check out the [plugin docs](docs/client-plugins.md) on how to create and compile plugins for **Deteccted**.

## 🏃 Running the client

## Using 🐳 Docker

Build the Deteccted image:

```bash
docker build --build-arg PLUGIN_DIR=/path/to/plugins --target=app -t deteccted . 
```

Run the Deteccted container:

```bash
docker run -v ./config.yaml:/deteccted/src/config.yaml deteccted 
```

## Standalone

```bash
go build -o deteccted . 
./deteccted #--help for all the available flags
```