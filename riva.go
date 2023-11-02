package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func randomString(n int) (string, error) {
	const lettersAndDigits = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = lettersAndDigits[int(b[i])%len(lettersAndDigits)]
	}
	return string(b), nil
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	maxMemory := int64(1024 * 1024 * 200) // 200 MB
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error parsing multipart form: %s", err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error getting file: %s", err)
	}
	defer file.Close()

	// Random filename, + extension.
	filename, err := randomString(6)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error generating random string: %s", err)
	}

	filename = filename + header.Filename[len(header.Filename)-4:]

	out, err := os.Create("files/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", "localhost:8080/"+filename)
}

func main() {
	os.MkdirAll("files", 0777)
	http.HandleFunc("/", fileUploadHandler)
	http.ListenAndServe(":8080", nil)
}
