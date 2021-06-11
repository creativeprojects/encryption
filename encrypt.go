package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	N = 32768
	r = 8
	p = 1
)

type Encrypter struct {
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func NewEncrypter(passphrase, salt []byte) (*Encrypter, error) {
	key, err := scrypt.Key(passphrase, salt, N, r, p, 32)
	if err != nil {
		return nil, fmt.Errorf("cannot generate an encryption key for AES: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cannot use encryption key for AES: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot create a Galois Counter Mode for AES: %w", err)
	}
	return &Encrypter{
		key:    key,
		aesgcm: aesgcm,
		nonce:  []byte{},
	}, nil
}

func (e *Encrypter) Key() []byte {
	return e.key
}

func (e *Encrypter) Nonce() []byte {
	if len(e.nonce) > 0 {
		return e.nonce
	}
	e.nonce = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, e.nonce); err != nil {
		// there's not much we can do here
		panic(err)
	}
	return e.nonce
}

func (e *Encrypter) Encrypt(plaintext []byte) []byte {
	ciphertext := e.aesgcm.Seal(nil, e.Nonce(), plaintext, nil)
	return ciphertext
}
