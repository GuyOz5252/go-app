package websocket

import (
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	userId     string
	hub        *Hub
	connection *websocket.Conn
	send       chan *WSMessage
}

func NewClient(userId string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		userId:     userId,
		hub:        hub,
		connection: conn,
		send:       make(chan *WSMessage, 10),
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
			c.hub.sendWSError(c.userId, err)
			break
		}
		var wsMessage WSMessage
		if err := json.Unmarshal(messageBytes, &wsMessage); err != nil {
			c.hub.sendWSError(c.userId, err)
		} else {
			ctx := context.Background()

			switch wsMessage.Type {
			case NewMessage:
				c.hub.deliverMessage(ctx, &wsMessage)
			case MessageUserAck:
				c.hub.sendWSMessage(wsMessage.DestinationUserId, &wsMessage)
			case MessageUserReadAck:
				c.hub.sendWSMessage(wsMessage.DestinationUserId, &wsMessage)
			case UserTypingStart:
				c.hub.sendWSMessageToDestinationChat(ctx, &wsMessage)
			case UserTypingEnd:
				c.hub.sendWSMessageToDestinationChat(ctx, &wsMessage)
			default:
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
