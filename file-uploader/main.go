package main

import (
    "file-uploader/server"
    "file-uploader/utils"
)

func main() {
    utils.Info("Starting File Uploader Server on :8080")
    server.StartServer("8080")
}
