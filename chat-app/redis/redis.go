package redis

import (
    "context"
    "chat-app/config"
    "chat-app/utils"
    "github.com/go-redis/redis/v8"
)

var (
    Rdb *redis.Client
    Ctx = context.Background()
)

func Init() {
    Rdb = redis.NewClient(&redis.Options{
        Addr: config.RedisAddr,
        DB:   config.RedisDB,
    })
    _, err := Rdb.Ping(Ctx).Result()
    if err != nil {
        utils.Error("Failed to connect to Redis: " + err.Error())
    } else {
        utils.Info("Connected to Redis")
    }
}

func Publish(channel, message string) {
    Rdb.Publish(Ctx, channel, message)
}

func Subscribe(channel string) *redis.PubSub {
    return Rdb.Subscribe(Ctx, channel)
}
