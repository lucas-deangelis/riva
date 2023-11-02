package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rivaupload <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if err := UploadFile(filename); err != nil {
		fmt.Println("Error uploading file:", err)
		os.Exit(1)
	}

	fmt.Println("File uploaded successfully.")
}

func UploadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	writer.Close()

	targetURL := "http://localhost:8080/"
	req, err := http.NewRequest("POST", targetURL, &buffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Response:", string(responseBody))

	return nil
}
