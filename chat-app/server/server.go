package server

import (
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var rooms = make(map[string]*Room)

func GetRoom(name string) *Room {
    if room, ok := rooms[name]; ok {
        return room
    }
    room := NewRoom(name)
    rooms[name] = room
    return room
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
    roomName := r.URL.Query().Get("room")
    if roomName == "" {
        roomName = "default"
    }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    client := &Client{
        Conn: conn,
        Room: GetRoom(roomName),
        Send: make(chan []byte, 256),
    }

    client.Room.Register <- client

    go client.WritePump()
    client.ReadPump()
}
