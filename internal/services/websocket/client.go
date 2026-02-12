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

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.connection.Close()
	}()
	for {
		_, _, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// TODO: Log
			}
			break
		}
		
		// TODO: Handle messages
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
