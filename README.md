# Chat App

A scalable WebSocket chat application in Go with multiple rooms and Redis pub/sub.

## Run

1. Make sure Redis is running on localhost:6379
2. Run the server:

```bash
go run main.go
```

3. Connect using a WebSocket client:

```
ws://localhost:8080/ws?room=room_name
```
