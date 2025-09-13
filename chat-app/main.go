package main

import (
    "chat-app/redis"
    "chat-app/server"
    "chat-app/utils"
    "net/http"
)

func main() {
    redis.Init()
    http.HandleFunc("/ws", server.ServeWs)

    utils.Info("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}
