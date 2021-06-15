package main

import (
	"net/http"
)

type ResponseWriter struct {
	encrypter *Encrypter
	nonce     Nonce
	writer    http.ResponseWriter
}

func NewResponseWriter(encrypter *Encrypter, rw http.ResponseWriter) *ResponseWriter {
	nonce := NewNonce(encrypter.NonceSize())

	header := rw.Header()
	contentType := header.Get(ContentTypeHeader)
	if contentType != "" {
		header.Set(EmbeddedContentTypeHeader, contentType)
	}
	header.Set(ContentTypeHeader, ContentType)

	return &ResponseWriter{
		encrypter: encrypter,
		nonce:     nonce,
		writer:    rw,
	}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.writer.Header()
}

func (rw *ResponseWriter) Write(source []byte) (int, error) {
	data := rw.encrypter.EncryptWithNonce(source, rw.nonce)
	_, err := rw.writer.Write(data)
	// the http server is expecting the original length written back
	return len(source), err
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.writer.WriteHeader(statusCode)
}

var _ http.ResponseWriter = &ResponseWriter{}
