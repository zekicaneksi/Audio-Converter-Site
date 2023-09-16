package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type test struct {
	Abc string
}

func setupServer() *http.ServeMux {
	server := http.NewServeMux()

	server.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		// Get handler for filename, size and headers
		file, handler, err := r.FormFile("audioFile")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		if handler.Size > 30000000 {
			http.Error(w, "400 file size too big.", http.StatusBadRequest)
			return
		}

		// Save the file locally

		// Create file
		dst, err := os.Create("uploaded_files/" + handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			defer dst.Close()
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// response
		w.Header().Set("Content-Type", "application/json")

		abc := test{
			Abc: "123",
		}

		json.NewEncoder(w).Encode(abc)
	})

	return server
}

func main() {

	server := setupServer()

	log.Fatal(http.ListenAndServe(":8080", server))
}
