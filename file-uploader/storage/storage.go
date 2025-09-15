package storage

import (
    "file-uploader/utils"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
)

func SaveFile(file multipart.File, filename string, uploadPath string) error {
    os.MkdirAll(uploadPath, os.ModePerm)
    out, err := os.Create(filepath.Join(uploadPath, filename))
    if err != nil {
        utils.Error("Failed to create file: " + err.Error())
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        utils.Error("Failed to save file: " + err.Error())
        return err
    }
    return nil
}
