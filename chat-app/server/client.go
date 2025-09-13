package server

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    Conn *websocket.Conn
    Room *Room
    Send chan []byte
}

func (c *Client) ReadPump() {
    defer func() {
        c.Room.Unregister <- c
        c.Conn.Close()
    }()

    for {
        _, msg, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }
        c.Room.Broadcast <- msg
    }
}

func (c *Client) WritePump() {
    defer c.Conn.Close()
    for msg := range c.Send {
        c.Conn.WriteMessage(websocket.TextMessage, msg)
    }
}
