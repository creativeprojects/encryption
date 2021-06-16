package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

type Decrypter struct {
	key    []byte
	aesgcm cipher.AEAD
}

func NewDecrypter(passphrase, salt []byte) (*Decrypter, error) {
	key, err := scrypt.Key(passphrase, salt, scrypt_N, scrypt_r, scrypt_p, 32)
	if err != nil {
		return nil, fmt.Errorf("cannot generate encryption key: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cannot use encryption key for AES: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot create a Galois Counter Mode for AES: %w", err)
	}
	return &Decrypter{
		key:    key,
		aesgcm: aesgcm,
	}, nil
}

func (d *Decrypter) Decrypt(ciphertext []byte, nonce Nonce) ([]byte, error) {
	if len(nonce) != d.aesgcm.NonceSize() {
		return nil, errors.New("incorrect nonce length given to GCM")
	}
	plaintext, err := d.aesgcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}
