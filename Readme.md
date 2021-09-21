# Deteccted

## About

Deteccted is the client for **detecctor**. It should be installed on the same server where the _service nodes_ run. It will
periodically ping the server with the status of the service node. The client also supports plugins which are bound to
the plugins and commands on the server.

## Configuration

```yaml
client:
  server:
    host: localhost
    port: 7777
pingInterval: 60
pluginDir: "/usr/deteccted/plugins"
plugins:
  - "examplePlugin"
```

## Running the client

## Docker

```bash
todo  
```

## Standalone

```bash
go build . 
./main #--help
```