package encryption

import "errors"

var (
	ErrorNotEncrypted            = errors.New("content is not encrypted")
	ErrorMissingEncryptionHeader = errors.New("missing encryption header")
	ErrorInvalidEncryptionHeader = errors.New("invalid encryption header")
	ErrorInvalidEncryptionKey    = errors.New("invalid encryption key")
)
