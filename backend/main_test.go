package main

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func makeRequestWithFile(t *testing.T, inputFileName string, outputFileName string, portNumber string, readCount int) *http.Response {
	body, writer := io.Pipe()
	mwriter := multipart.NewWriter(writer)

	go func() {
		defer writer.Close()
		defer mwriter.Close()

		w, err := mwriter.CreateFormFile("audioFile", "test_files/"+outputFileName)
		if err != nil {
			t.Log("Error creating form file")
			t.Log(err)
			return
		}

		for i := 0; i < readCount; i++ {
			in, err := os.Open("test_files/" + inputFileName)
			if err != nil {
				t.Log("Error opening file")
				t.Log(err)
				return
			}
			defer in.Close()

			if _, err := io.Copy(w, in); err != nil {
				t.Log("Error copying file")
				t.Log(err)
				return
			}
		}

		if err := mwriter.Close(); err != nil {
			t.Log("Error form data")
			t.Log(err)
			return
		}
	}()

	res, err := http.Post("http://localhost:"+portNumber+"/upload", mwriter.FormDataContentType(), body)
	if err != nil {
		t.Fatal("could not make the request")
		return nil
	}

	return res
}

func TestValidFileUpload(t *testing.T) {

	portNumber := "8080"

	server := setupServer()
	go http.ListenAndServe(":"+portNumber, server)

	res := makeRequestWithFile(t, "valid_file.mp3", "valid.mp3", portNumber, 1)

	var resValue uploadResponse
	err := json.NewDecoder(res.Body).Decode(&resValue)
	if err != nil {
		t.Fatal("error decoding response")
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal("Failed should have uploaded")
	}

}

func TestBigFileUpload(t *testing.T) {

	portNumber := "8081"

	server := setupServer()
	go http.ListenAndServe(":"+portNumber, server)

	res := makeRequestWithFile(t, "valid_file.mp3", "big.mp3", portNumber, 15)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Failed, should have not uploaded")
	}

}

func TestInvalidFileUpload(t *testing.T) {

	portNumber := "8082"

	server := setupServer()
	go http.ListenAndServe(":"+portNumber, server)

	res := makeRequestWithFile(t, "invalid_file.pdf", "invalid.pdf", portNumber, 1)

	if res.StatusCode != http.StatusNotAcceptable {
		t.Fatal("Failed should not have accepted")
	}

}
