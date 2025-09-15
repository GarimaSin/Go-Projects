# File Uploader

A scalable file uploader in Go.

## Run

1. Run the server:

```bash
go run main.go
```

2. Upload a file using curl:

```bash
curl -F "file=@/path/to/file" http://localhost:8080/upload
```

3. Files will be saved in the `uploads/` directory.
