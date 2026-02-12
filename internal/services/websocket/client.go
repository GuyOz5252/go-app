package websocket

import (
	"encoding/json"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/gorilla/websocket"
)

type Client struct {
	userId          string
	hub             *Hub
	connection      *websocket.Conn
	send            chan *core.WSMessage
	handlerResolver map[string]func(m *core.WSMessage)
}

func NewClient(userId string, hub *Hub, conn *websocket.Conn, hr map[string]func(m *core.WSMessage)) *Client {
	return &Client{
		userId:          userId,
		hub:             hub,
		connection:      conn,
		send:            make(chan *core.WSMessage, 10),
		handlerResolver: hr,
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		c.hub.Unregister <- c
		c.connection.Close()
	}()
	for {
		_, messageBytes, err := c.connection.ReadMessage()
		if err != nil {
			// TODO: Log
			break
		}
		var wsMessage core.WSMessage
		if err := json.Unmarshal(messageBytes, &wsMessage); err != nil {
			// TODO: Log
		} else {
			c.handlerResolver[wsMessage.Type](&wsMessage)
		}
	}
}

func (c *Client) WriteMessages() {
	defer c.connection.Close()

	for webSocketMessage := range c.send {
		if err := c.connection.WriteJSON(webSocketMessage); err != nil {
			return
		}
	}

	c.connection.WriteMessage(websocket.CloseMessage, []byte{})
}
