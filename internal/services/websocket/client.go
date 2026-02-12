package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub        *Hub
	connection *websocket.Conn
	send       chan []byte
	userId     string
}

func NewClient(hub *Hub, conn *websocket.Conn, userId string) *Client {
	return &Client{
		hub:        hub,
		connection: conn,
		send:       make(chan []byte, 256),
		userId:     userId,
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		c.hub.Unregister <- c
		c.connection.Close()
	}()
	for {
		_, _, err := c.connection.ReadMessage()
		if err != nil {
			// TODO: Log
			break
		}

		// TODO: Handle messages
	}
}

func (c *Client) WriteMessages() {
	defer c.connection.Close()

	for message := range c.send {
		if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}

	c.connection.WriteMessage(websocket.CloseMessage, []byte{})
}
