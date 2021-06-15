package main

import (
	"fmt"
	"net/http"
)

type EncryptionHandler struct {
	encrypter *Encrypter
}

func NewEncryptionHandler(passphrase []byte) (*EncryptionHandler, error) {
	encrypter, err := NewEncrypter(passphrase, []byte(DefaultSalt))
	if err != nil {
		return nil, fmt.Errorf("cannot create encrypter: %w", err)
	}
	return &EncryptionHandler{
		encrypter: encrypter,
	}, nil
}

func (h *EncryptionHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// writer := NewResponseWriter(h.encrypter, rw)
}

var _ http.Handler = &EncryptionHandler{}
