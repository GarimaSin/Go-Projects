package server

import (
    "chat-app/redis"
    "chat-app/utils"
)

type Room struct {
    Name       string
    Clients    map[*Client]bool
    Broadcast  chan []byte
    Register   chan *Client
    Unregister chan *Client
}

func NewRoom(name string) *Room {
    room := &Room{
        Name:       name,
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan []byte),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
    }

    go room.run()
    go room.redisSubscribe()

    return room
}

func (r *Room) run() {
    for {
        select {
        case client := <-r.Register:
            r.Clients[client] = true
        case client := <-r.Unregister:
            if _, ok := r.Clients[client]; ok {
                delete(r.Clients, client)
                close(client.Send)
            }
        case msg := <-r.Broadcast:
            redis.Publish(r.Name, string(msg))
            for client := range r.Clients {
                select {
                case client.Send <- msg:
                default:
                    close(client.Send)
                    delete(r.Clients, client)
                }
            }
        }
    }
}

func (r *Room) redisSubscribe() {
    sub := redis.Subscribe(r.Name)
    ch := sub.Channel()
    for msg := range ch {
        for client := range r.Clients {
            select {
            case client.Send <- []byte(msg.Payload):
            default:
                close(client.Send)
                delete(r.Clients, client)
            }
        }
    }
}
