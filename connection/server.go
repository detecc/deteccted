package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"github.com/detecc/detecctor/shared"
	"github.com/detecc/deteccted/cache"
	"github.com/detecc/deteccted/config"
	"github.com/detecc/deteccted/plugin"
	"sync"
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

	plugin.LoadPlugins()
	// send the auth request
	once := sync.Once{}
	once.Do(func() {
		client.sendMessage(&shared.Payload{
			Id:             "1",
			ServiceNodeKey: "a7309be127cf7g9127309",
			Data:           config.GetClientConfiguration().Client.AuthPassword,
			Command:        "/auth",
			Success:        true,
			Error:          "",
		})
	})
}

func (c *Client) listenForIncomingMessages() {
	defer c.conn.Close()
	for {
		message, _ := bufio.NewReader(c.conn).ReadString('\n')
		if message != "" {
			c.handleMessage(message)
		}
	}
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
