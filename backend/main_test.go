package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func makeRequestWithFile(t *testing.T, inputFileName string, outputFileName string, portNumber string, readCount int) bool {
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
	if err != nil || res.StatusCode != 200 {
		t.Log(err)
		return false
	}
	return true
}

func TestValidFileUpload(t *testing.T) {

	portNumber := "8080"

	server := setupServer()
	go http.ListenAndServe(":"+portNumber, server)

	err := makeRequestWithFile(t, "valid_file.mp3", "valid.mp3", portNumber, 1)
	if !err {
		t.Fatal("Failed should have uploaded")
	}

}

func TestBigFileUpload(t *testing.T) {

	portNumber := "8081"

	server := setupServer()
	go http.ListenAndServe(":"+portNumber, server)

	err := makeRequestWithFile(t, "valid_file.mp3", "big.mp3", portNumber, 15)
	if err {
		t.Fatal("Failed, should have not uploaded")
	}

}
