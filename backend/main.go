package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type uploadResponse struct {
	Response string
}

func setupServer() *http.ServeMux {

	server := http.NewServeMux()
	
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("hello"))
	})

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

		// Run ffprobe on the file
		ffprobe_res, err := runffprobe(&handler.Filename)

		// Delete the file
		if err_delete := os.Remove("uploaded_files/" + handler.Filename); err_delete != nil {
			log.Println("could not remove the file")
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		// response
		w.Header().Set("Content-Type", "application/json")

		res := uploadResponse{
			Response: *ffprobe_res,
		}

		json.NewEncoder(w).Encode(res)
	})

	return server
}

func runffprobe(filename *string) (*string, error) {

	out, err := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "uploaded_files/"+*filename).Output()
	if err != nil {
		return nil, errors.New("ffprobe could not inspect or find the file, please provide a valid file")
	}

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(out), &jsonMap)

	b, err := json.Marshal(jsonMap["format"])
	if err != nil {
		return nil, errors.New("error decoding ffprobe output")
	}
	result := string(b)

	return &result, nil

}

func main() {

	server := setupServer()

	log.Fatal(http.ListenAndServe(":8080", server))
}
