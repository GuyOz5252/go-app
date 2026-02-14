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
	handlerResolver map[core.WSMessageType]func(m *core.WSMessage)
}

func NewClient(userId string, hub *Hub, conn *websocket.Conn, hr map[core.WSMessageType]func(m *core.WSMessage)) *Client {
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
			handlerFunc, ok := c.handlerResolver[wsMessage.Type]
			if (ok) {
				handlerFunc(&wsMessage)
			}
		}
	}
}

func (c *Client) WriteMessages() {
	defer c.connection.Close()

	for wsMessage := range c.send {
		if err := c.connection.WriteJSON(wsMessage); err != nil {
			return
		}
	}

	c.connection.WriteMessage(websocket.CloseMessage, []byte{})
}
