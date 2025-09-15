package server

import (
    "file-uploader/config"
    "file-uploader/storage"
    "file-uploader/utils"
    "fmt"
    "net/http"
    "time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        utils.Error("Failed to parse form: " + err.Error())
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        utils.Error("Failed to get file: " + err.Error())
        http.Error(w, "File missing", http.StatusBadRequest)
        return
    }
    defer file.Close()

    filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
    err = storage.SaveFile(file, filename, config.UploadPath)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", filename)))
}
