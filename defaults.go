package main

const (
	N                         = 32768
	r                         = 8
	p                         = 1
	ContentType               = "application/octet-stream"
	ContentTypeHeader         = "Content-Type"
	EmbeddedContentTypeHeader = "X-Content-Type"
)

var DefaultSalt = "Poor's Man Encrypted HTTP Transfer"
