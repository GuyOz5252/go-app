package websocket

import (
	"context"
	"encoding/json"

	"github.com/GuyOz5252/go-app/internal/core"
	"github.com/gorilla/websocket"
)

type Client struct {
	userId     string
	hub        *Hub
	connection *websocket.Conn
	send       chan *core.WSMessage
}

func NewClient(userId string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		userId:     userId,
		hub:        hub,
		connection: conn,
		send:       make(chan *core.WSMessage, 10),
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
			c.hub.sendWSError("", c.userId, err)
			break
		}
		var wsMessage core.WSMessage
		if err := json.Unmarshal(messageBytes, &wsMessage); err != nil {
			c.hub.sendWSError("", c.userId, err)
		} else {
			ctx := context.Background()

			switch wsMessage.Type {
			case core.NewMessage:
				c.hub.deliverMessage(ctx, wsMessage)
			case core.MessageUserAck:
				c.hub.deliverUserAcks(wsMessage)
			case core.MessageUserReadAck:
				c.hub.deliverUserAcks(wsMessage)
			case core.UserTypingStart:
				c.hub.deliverTyping(ctx, wsMessage)
			case core.UserTypingEnd:
				c.hub.deliverTyping(ctx, wsMessage)
			default:
			}
		}
	}
}

func (c *Client) WriteMessages() {
	defer c.connection.Close()

	for wsMessage := range c.send {
		if err := c.connection.WriteJSON(wsMessage); err != nil {
			break
		}
	}

	c.connection.WriteMessage(websocket.CloseMessage, []byte{})
}
