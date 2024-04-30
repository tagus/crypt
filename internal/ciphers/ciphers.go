package ciphers

type Cipher interface {
	Encrypt(string) ([]byte, error)
	Decrypt([]byte) (string, error)
}
