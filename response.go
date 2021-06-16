package encryption

import (
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

func AddAcceptEncryptionHeader(header http.Header) {
	header.Add(acceptEncodingHeader, ContentEncoding)
}

func DecryptResponse(resp *http.Response, passphrase []byte) ([]byte, error) {
	if resp.Header.Get(contentEncodingHeader) != ContentEncoding {
		return nil, ErrorNotEncrypted
	}
	header := strings.Split(resp.Header.Get(ContentEncryptionHeader), "|")
	if len(header) != 2 {
		return nil, ErrorMissingEncryptionHeader
	}

	salt, err := hex.DecodeString(header[0])
	if err != nil {
		return nil, ErrorInvalidEncryptionHeader
	}

	nonce, err := hex.DecodeString(header[1])
	if err != nil {
		return nil, ErrorInvalidEncryptionHeader
	}

	decrypter, err := NewDecrypter(passphrase, salt)
	if err != nil {
		return nil, ErrorInvalidEncryptionKey
	}

	ciphertext, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	plaintext, err := decrypter.Decrypt(ciphertext, nonce)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
