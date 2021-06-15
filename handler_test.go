package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	echoHandler *http.ServeMux
	passphrase  = []byte("test passphrase")
)

func init() {
	encrypter, err := NewEncrypter(passphrase, []byte(DefaultSalt))
	if err != nil {
		panic(err)
	}

	echoHandler = http.NewServeMux()
	echoHandler.HandleFunc("/plain", func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			rw.Header().Set("Content-Type", contentType)
		}

		_, err := io.Copy(rw, r.Body)
		if err != nil {
			panic(err)
		}
	})
	echoHandler.HandleFunc("/encrypted", func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			rw.Header().Set("Content-Type", contentType)
		}

		wrapper := NewResponseWriter(encrypter, rw)

		_, err := io.Copy(wrapper, r.Body)
		if err != nil {
			panic(err)
		}
	})
}

func TestWithNoEncryption(t *testing.T) {
	message := []byte("This is a test!")
	reader := bytes.NewReader(message)
	server := httptest.NewServer(echoHandler)
	client := server.Client()

	resp, err := client.Post(server.URL+"/plain", "text/plain", reader)
	require.NoError(t, err)

	defer resp.Body.Close()
	returned, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, message, returned)

	server.Close()
}

func TestWithEncryption(t *testing.T) {
	message := []byte("This is a test!")
	reader := bytes.NewReader(message)
	server := httptest.NewServer(echoHandler)
	client := server.Client()

	resp, err := client.Post(server.URL+"/encrypted", "text/plain", reader)
	require.NoError(t, err)

	defer resp.Body.Close()
	returned, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.NotEqual(t, message, returned)

	assert.Equal(t, "application/octet-stream", resp.Header.Get("Content-Type"))
	assert.Equal(t, "text/plain", resp.Header.Get("X-Content-Type"))

	server.Close()
}
