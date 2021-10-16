package connection

import (
	"bufio"
	"fmt"
	"github.com/detecc/deteccted/config"
	"github.com/detecc/deteccted/plugin"
	"github.com/detecc/detecctor/shared"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var client *Client
var once = sync.Once{}

type Client struct {
	conn net.Conn
}

// Start connects to the tcp server, and listens for incoming commands
func Start() {
	conf := config.GetClientConfiguration()

	plugin.GetPluginManager().LoadPlugins()

	once.Do(func() {
		// construct the Client only once
		client = &Client{
			conn: dial(conf.Client.Host, conf.Client.Port),
		}

		go client.listenForIncomingMessages()

		// send the auth request
		payload := shared.NewPayload(
			shared.ForClient(conf.ServiceNodeIdentifier),
			shared.WithData(conf.Client.AuthPassword),
			shared.ForCommand("/auth"),
			shared.Successful(),
		)
		client.sendMessage(&payload)
	})
}

// SendPayloadToServer sends a payload to the server if the client exists and is connected.
//It is just an exposure to the plugins that want to periodically transmit data without prior request.
func SendPayloadToServer(payload shared.Payload) error {
	if client != nil {
		return client.sendMessage(&payload)
	}
	return fmt.Errorf("client not initialized")
}

// SendToServer sends a payload to the server if the client exists and is connected.
//It is just an exposure to the plugins that want to periodically transmit data without prior request.
func SendToServer(opts ...shared.PayloadOption) error {
	if client != nil {
		payload := shared.NewPayload(opts...)
		return client.sendMessage(&payload)
	}
	return fmt.Errorf("client not initialized")
}

func (c *Client) listenForIncomingMessages() {
	conf := config.GetClientConfiguration().Client
	defer c.conn.Close()
	for {
		message, err := bufio.NewReader(c.conn).ReadString('\n')
		if err == io.EOF {
			//try to reconnect
			log.Println("Connection is down, reconnecting...")
			conn, reconnectErr := redial(conf.Host, conf.Port)
			if reconnectErr == nil {
				log.Println("Reconnected")
				c.conn = conn
				continue
			}

			// timeout for 30 seconds before retrying
			time.Sleep(30 * time.Second)
			continue
		}
		if message != "" {
			c.handleMessage(message)
		}
	}
}

// handleMessage processes the message and executes corresponding plugin
func (c *Client) handleMessage(message string) {
	log.Println("Received message from server:", message)
	// decode the payload
	payload := shared.NewPayload()
	err := shared.DecodePayload([]byte(message), &payload)
	if err != nil {
		log.Println("couldn't decode the payload:", err)
		return
	}

	switch payload.Command {
	case "/ping":
		payload.Data = "pong"
		c.sendMessage(&payload)
		break
	default:
		c.executePlugin(&payload)
		break
	}
}

// executePlugin executes the plugin.
func (c *Client) executePlugin(payload *shared.Payload) {
	// get the plugin
	mPlugin, err := plugin.GetPluginManager().GetPlugin(payload.Command)
	if err != nil {
		log.Println("command/plugin not found:", payload.Command)
		payload.SetError(fmt.Errorf("client handler for command %s not found", payload.Command))
		c.sendMessage(payload)
		return
	}

	switch mPlugin.GetMetadata().Type {
	case plugin.PluginTypeClientServer:
		// execute the plugin
		response, pluginErr := mPlugin.Execute(payload.Data)
		if pluginErr != nil {
			log.Println("Plugin execution returned error:", pluginErr)
			payload.SetError(pluginErr)
			c.sendMessage(payload)
			break
		}

		payload.Data = response
		payload.Success = true
		c.sendMessage(payload)
		break
	case plugin.PluginTypeClientOnly:
		payload.SetError(fmt.Errorf("plugin cannot handle the command %s", payload.Command))
		c.sendMessage(payload)
		break
	default:
		log.Println("unsupported plugin type:", mPlugin.GetMetadata().Type)
		payload.SetError(fmt.Errorf("unsupported plugin type:%s", mPlugin.GetMetadata().Type))
		c.sendMessage(payload)
	}
}

// sendMessage prepares a payload, packs it and sends it to the TCP server.
func (c *Client) sendMessage(payload *shared.Payload) error {
	log.Println("Sending payload", payload)
	message, err := shared.EncodePayload(payload)
	if err != nil {
		log.Println("cannot encode the payload:", err)
		return err
	}

	_, err = c.conn.Write([]byte(message))
	if err != nil {
		log.Println("couldn't send the payload:", err)
	}
	return err
}

// dial is used the first time when connecting to a server. If the connection fails immediately, exit.
func dial(host string, port int) net.Conn {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		log.Fatal(e)
	}
	return conn
}

// redial is used for reconnecting to the server after the client connection went down.
func redial(host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	return net.Dial("tcp", addr)
}
