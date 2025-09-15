package server

import "net/http"

func StartServer(port string) {
    http.HandleFunc("/upload", UploadHandler)
    http.ListenAndServe(":"+port, nil)
}
