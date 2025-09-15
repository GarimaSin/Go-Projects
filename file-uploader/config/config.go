package config

import "os"

var (
    UploadPath = getEnv("UPLOAD_PATH", "./uploads")
)

func getEnv(key, defaultVal string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultVal
}
