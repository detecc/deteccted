package connection

import (
	"bufio"
	"fmt"
	"github.com/detecc/deteccted/cache"
	"github.com/detecc/deteccted/config"
	"github.com/detecc/deteccted/plugin"
	"github.com/detecc/detecctor/shared"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
}

// Start connects to the tcp server, and listens for incoming commands
func Start() {
	client := &Client{
		conn: dial("localhost", 7777),
	}
	go client.listenForIncomingMessages()

	plugin.GetPluginManager().LoadPlugins()

	conf := config.GetClientConfiguration()
	// send the auth request
	once := sync.Once{}
	once.Do(func() {
		client.sendMessage(&shared.Payload{
			Id:             "1",
			ServiceNodeKey: conf.ServiceNodeIdentifier,
			Data:           conf.Client.AuthPassword,
			Command:        "/auth",
			Success:        true,
			Error:          "",
		})
	})
}

func (c *Client) listenForIncomingMessages() {
	conf := config.GetClientConfiguration().Client
	defer c.conn.Close()
	for {
		if !c.isConnectionAlive() {
			//try to reconnect
			log.Println("Connection is down, reconnecting...")
			c.conn = dial(conf.Host, conf.Port)
			time.Sleep(1 * time.Second)
			continue
		}
		message, _ := bufio.NewReader(c.conn).ReadString('\n')
		if message != "" {
			c.handleMessage(message)
		}
	}
}
// isConnectionAlive Checks if the connection is alive.
func (c *Client) isConnectionAlive() bool {
	one := make([]byte, 1)
	err := c.conn.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	if err != nil {
		return false
	}
	if _, err := c.conn.Read(one); err == io.EOF {
		return false
	}
	return true
}

// handleMessage processes the message and executes corresponding plugin
func (c *Client) handleMessage(message string) {
	log.Println("Received message from server:", message)
	// decode the payload
	payload := shared.Payload{}
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
		command, ok := cache.Memory().Get(payload.Command)
		if !ok {
			log.Println("command not found")
			payload.Error = fmt.Sprintf("Client handler for command %s not found!", payload.Command)
			payload.Success = false
			c.sendMessage(&payload)
			return
		}

		response, err := command.(plugin.Handler).Execute(payload.Data)
		if err != nil {
			log.Println("Plugin execution returned error:", err)
			payload.Error = err.Error()
			payload.Success = false
			c.sendMessage(&payload)
			return
		}

		payload.Data = response

		c.sendMessage(&payload)
		break
	}
}

// sendMessage prepares a payload, packs it and sends it to the tcp server
func (c *Client) sendMessage(payload *shared.Payload) {
	log.Println("Sending payload", payload)
	message, err := shared.EncodePayload(payload)
	if err != nil {
		log.Println("cannot encode the payload:", err)
		return
	}

	_, err = c.conn.Write([]byte(message))
	if err != nil {
		log.Println("couldn't send the payload:", err)
	}
}

func dial(host string, port int) net.Conn {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		log.Fatal(e)
	}
	return conn
}
