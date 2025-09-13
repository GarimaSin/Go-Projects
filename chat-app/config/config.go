package config

import "os"

var (
    RedisAddr = getEnv("REDIS_ADDR", "localhost:6379")
    RedisDB   = 0
)

func getEnv(key, defaultVal string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultVal
}
