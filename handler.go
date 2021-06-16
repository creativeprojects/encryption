package encryption

import (
	"fmt"
	"net/http"
	"strings"
)

type Handler struct {
	encrypter *Encrypter
	next      http.Handler
}

func NewHandler(passphrase, salt []byte, nextHandler http.Handler) (*Handler, error) {
	encrypter, err := NewEncrypter(passphrase, []byte(salt))
	if err != nil {
		return nil, fmt.Errorf("cannot create encrypter: %w", err)
	}
	return &Handler{
		encrypter: encrypter,
		next:      nextHandler,
	}, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Header.Get(acceptEncodingHeader), ContentEncoding) {
		h.next.ServeHTTP(rw, req)
		return
	}
	writer := NewEncryptedWriter(h.encrypter, rw)
	h.next.ServeHTTP(writer, req)
}

var _ http.Handler = &Handler{}
