package encryption

import (
	"crypto/rand"
	"io"
)

type Nonce []byte

func NewNonce(size int) Nonce {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		// there's not much we can do here
		panic(err)
	}
	return nonce
}
