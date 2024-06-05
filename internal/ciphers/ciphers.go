package ciphers

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type Cipher interface {
	Encrypt(string) ([]byte, error)
	Decrypt([]byte) (string, error)
}
