package main

import (
	"log"
	"net/http"

	"github.com/creativeprojects/encryption"
)

const serverURL = "http://localhost:3001/"

func main() {
	log.Printf("Calling %s", serverURL)

	req, err := http.NewRequest(http.MethodGet, serverURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// add the accept-encoding header saying we accept our encrypted content
	encryption.AddAcceptEncryptionHeader(req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// decrypt the response
	passphrase := []byte("I have been eating too much chocolate")
	message, err := encryption.DecryptResponse(resp, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Answer: %q\n", string(message))
}
