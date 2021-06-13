package main

import "net/http"

type EncryptionHandler struct {
	passphrase []byte
	next       http.Handler
}

func NewEncryptionHandler(passphrase []byte, nextHandler http.Handler) *EncryptionHandler {
	return &EncryptionHandler{
		passphrase: passphrase,
		next:       nextHandler,
	}
}

func (h *EncryptionHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// pass-through
	h.next.ServeHTTP(rw, req)
}

var _ http.Handler = &EncryptionHandler{}
