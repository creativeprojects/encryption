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
)

func init() {
	echoHandler = http.NewServeMux()
	echoHandler.HandleFunc("/echo", func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		_, err := io.Copy(rw, r.Body)
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

	resp, err := client.Post(server.URL+"/echo", "text/plain", reader)
	require.NoError(t, err)

	defer resp.Body.Close()
	returned, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, message, returned)

	server.Close()
}
