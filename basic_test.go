package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptionKeyIsPermanentForGivenSalt(t *testing.T) {
	passphrase := "Too much chocolate"
	salt1 := "Way too salty"
	salt2 := "Not enough salt"

	encrypter1, err := NewEncrypter([]byte(passphrase), []byte(salt1))
	require.NoError(t, err)

	encrypter2, err := NewEncrypter([]byte(passphrase), []byte(salt1))
	require.NoError(t, err)

	encrypter3, err := NewEncrypter([]byte(passphrase), []byte(salt2))
	require.NoError(t, err)

	assert.Equal(t, encrypter1.Key(), encrypter2.Key())
	assert.NotEqual(t, encrypter1.Key(), encrypter3.Key())
}

func TestEncryptDecrypt(t *testing.T) {
	passphrase := "Too much chocolate"
	salt := "Way too salty"

	encrypter, err := NewEncrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)

	message := `    inet 127.0.0.1/8 scope host lo
    inet 192.168.213.18/24 brd 192.168.213.255 scope global eno2
    inet 192.168.211.18/24 brd 192.168.211.255 scope global br0
    inet 192.168.211.80/24 brd 192.168.211.255 scope global secondary br0
    inet 192.168.211.82/24 brd 192.168.211.255 scope global secondary br0
    inet 172.19.0.1/16 brd 172.19.255.255 scope global br-358e25c4c8ba
    inet 172.22.0.1/16 brd 172.22.255.255 scope global br-4064b1aa9e42
    inet 172.18.0.1/16 brd 172.18.255.255 scope global br-7530e767a981
    inet 172.23.0.1/16 brd 172.23.255.255 scope global br-868b07e0737c
    inet 172.24.0.1/16 brd 172.24.255.255 scope global br-e52b418d6cdd
    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
    inet 192.168.22.2/32 scope global wg0
`
	ciphertext, nonce := encrypter.Encrypt([]byte(message))

	decrypter, err := NewDecrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)
	plaintext, err := decrypter.Decrypt(ciphertext, nonce)
	require.NoError(t, err)

	assert.Equal(t, message, string(plaintext))
}

func TestEncryptDecryptEmptyMessage(t *testing.T) {
	passphrase := "Too much chocolate"
	salt := "Way too salty"

	encrypter, err := NewEncrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)

	message := ""
	ciphertext, nonce := encrypter.Encrypt([]byte(message))

	decrypter, err := NewDecrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)
	plaintext, err := decrypter.Decrypt(ciphertext, nonce)
	require.NoError(t, err)

	assert.Equal(t, message, string(plaintext))
}

func TestDecryptIncorrectNonce(t *testing.T) {
	passphrase := "Too much chocolate"
	salt := "Way too salty"

	encrypter, err := NewEncrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)

	message := ""
	ciphertext, _ := encrypter.Encrypt([]byte(message))
	_, nonce := encrypter.Encrypt([]byte(message))

	decrypter, err := NewDecrypter([]byte(passphrase), []byte(salt))
	require.NoError(t, err)
	_, err = decrypter.Decrypt(ciphertext, nonce)
	require.Error(t, err)
	t.Log(err)
}

func TestDecryptIncorrectSalt(t *testing.T) {
	passphrase := "Too much chocolate"
	salt1 := "Way too salty"
	salt2 := "Not enough salt"

	encrypter, err := NewEncrypter([]byte(passphrase), []byte(salt1))
	require.NoError(t, err)

	message := ""
	ciphertext, nonce := encrypter.Encrypt([]byte(message))

	decrypter, err := NewDecrypter([]byte(passphrase), []byte(salt2))
	require.NoError(t, err)
	_, err = decrypter.Decrypt(ciphertext, nonce)
	require.Error(t, err)
	t.Log(err)
}
