package encryption

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testHandler http.Handler
	passphrase  = []byte("test passphrase")
	testSalt    = []byte("not very salty")
)

func init() {
	var err error

	echoHandler := http.NewServeMux()
	echoHandler.HandleFunc("/echo", func(rw http.ResponseWriter, r *http.Request) {
		// copy content-type from request
		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			rw.Header().Set("Content-Type", contentType)
		}

		// copy body from request
		_, err := io.Copy(rw, r.Body)
		if err != nil {
			panic(err)
		}
	})

	testHandler, err = NewHandler(passphrase, testSalt, echoHandler)
	if err != nil {
		panic(err)
	}
}

func TestWithNoEncryption(t *testing.T) {
	message := []byte("This is a test!")
	reader := bytes.NewReader(message)
	server := httptest.NewServer(testHandler)
	client := server.Client()

	resp, err := client.Post(server.URL+"/echo", "text/plain", reader)
	require.NoError(t, err)
	defer resp.Body.Close()

	returned, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, "text/plain", resp.Header.Get(contentTypeHeader))
	assert.Equal(t, strconv.Itoa(len(message)), resp.Header.Get(contentLengthHeader))
	assert.Empty(t, resp.Header.Get(contentEncodingHeader))
	assert.Empty(t, resp.Header.Get(ContentEncryptionHeader))
	assert.Equal(t, message, returned)

	// t.Log(resp.Header)
	// t.Log(hex.Dump(returned))

	server.Close()
}

func TestWithEncryption(t *testing.T) {
	message := []byte("This is a test!")
	reader := bytes.NewReader(message)
	server := httptest.NewServer(testHandler)
	client := server.Client()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/echo", reader)
	require.NoError(t, err)
	req.Header.Set(acceptEncodingHeader, ContentEncoding)
	req.Header.Set(contentTypeHeader, "text/plain")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	returned, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, "text/plain", resp.Header.Get(contentTypeHeader))
	assert.Equal(t, strconv.Itoa(len(returned)), resp.Header.Get(contentLengthHeader))
	assert.Equal(t, ContentEncoding, resp.Header.Get(contentEncodingHeader))
	assert.NotEqual(t, message, returned)

	// t.Log(resp.Header)
	// t.Log(hex.Dump(returned))

	// check salt and nonce
	encryption := strings.Split(resp.Header.Get(ContentEncryptionHeader), "|")
	if len(encryption) != 2 {
		t.Errorf("invalid encryption header: %q", resp.Header.Get(ContentEncryptionHeader))
	}
	salt, err := hex.DecodeString(encryption[0])
	require.NoError(t, err)

	nonce, err := hex.DecodeString(encryption[1])
	require.NoError(t, err)

	decrypter, err := NewDecrypter(passphrase, salt)
	require.NoError(t, err)
	plaintext, err := decrypter.Decrypt(returned, nonce)
	require.NoError(t, err)

	// decryption successfull
	assert.Equal(t, message, plaintext)

	server.Close()
}
