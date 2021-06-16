package encryption

const (
	scrypt_N              = 32768
	scrypt_r              = 8
	scrypt_p              = 1
	acceptEncodingHeader  = "Accept-Encoding"
	contentLengthHeader   = "Content-Length"
	contentTypeHeader     = "Content-Type"
	contentEncodingHeader = "Content-Encoding"

	ContentEncoding         = "aesgcm"
	ContentEncryptionHeader = "X-Content-Encryption"
)
