package encryption

import (
	"fmt"
	"net/http"
	"sync"
)

type EncryptedWriter struct {
	encrypter  *Encrypter
	nonce      Nonce
	writer     http.ResponseWriter
	sendHeader sync.Once
}

func NewEncryptedWriter(encrypter *Encrypter, rw http.ResponseWriter) *EncryptedWriter {
	nonce := NewNonce(encrypter.NonceSize())

	return &EncryptedWriter{
		encrypter: encrypter,
		nonce:     nonce,
		writer:    rw,
	}
}

func (rw *EncryptedWriter) Header() http.Header {
	return rw.writer.Header()
}

func (rw *EncryptedWriter) Write(source []byte) (int, error) {
	data := rw.encrypter.EncryptWithNonce(source, rw.nonce)
	rw.writeHeader()
	_, err := rw.writer.Write(data)
	// the caller is expecting the original length back
	return len(source), err
}

func (rw *EncryptedWriter) WriteHeader(statusCode int) {
	rw.writeHeader()
	rw.writer.WriteHeader(statusCode)
}

func (rw *EncryptedWriter) writeHeader() {
	rw.sendHeader.Do(func() {
		header := rw.writer.Header()
		header.Add(contentEncodingHeader, ContentEncoding)
		header.Add("Vary", "Accept-Encoding")
		header.Set(ContentEncryptionHeader, fmt.Sprintf("%x|%x", rw.encrypter.Salt(), rw.nonce))
	})
}

// Flush implements http.Flusher
func (rw *EncryptedWriter) Flush() {
	if flusher, ok := rw.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

var (
	_ http.Flusher        = &EncryptedWriter{}
	_ http.ResponseWriter = &EncryptedWriter{}
)
