package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

type Encrypter struct {
	salt   []byte
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func NewEncrypter(passphrase, salt []byte) (*Encrypter, error) {
	key, err := scrypt.Key(passphrase, salt, scrypt_N, scrypt_r, scrypt_p, 32)
	if err != nil {
		return nil, fmt.Errorf("cannot generate an encryption key: %w", err)
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
		salt:   salt,
		key:    key,
		aesgcm: aesgcm,
		nonce:  []byte{},
	}, nil
}

func (e *Encrypter) NonceSize() int {
	return e.aesgcm.NonceSize()
}

func (e *Encrypter) Salt() []byte {
	return e.salt
}

func (e *Encrypter) Key() []byte {
	return e.key
}

func (e *Encrypter) Encrypt(plaintext []byte) ([]byte, Nonce) {
	nonce := NewNonce(e.aesgcm.NonceSize())
	ciphertext := e.EncryptWithNonce(plaintext, nonce)
	return ciphertext, nonce
}

func (e *Encrypter) EncryptWithNonce(plaintext []byte, nonce Nonce) []byte {
	ciphertext := e.aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext
}
