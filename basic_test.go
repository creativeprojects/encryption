package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	passphrase := "Too much chocolate"
	salt := "Way too salty"

	encrypter1, err := NewEncrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)

	encrypter2, err := NewEncrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)

	encrypter3, err := NewEncrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)

	t.Log(encrypter1.Key())
	assert.Equal(t, encrypter1.Key(), encrypter2.Key())
	assert.Equal(t, encrypter1.Key(), encrypter3.Key())
}
