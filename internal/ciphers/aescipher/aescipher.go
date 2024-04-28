package aescipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/tagus/crypt/internal/ciphers"
)

// AESCipher uses the AES encryption algorithm to encrypt and decrypt the provided data
// while ensuring that the data is signed with a hash of the provided password in order
// to ensure data integrity and provide password validation
type AESCipher struct {
	key   [32]byte
	block cipher.Block
}

func New(pwd string) (*AESCipher, error) {
	key := ciphers.ComputeHash(pwd)
	var keyBits [32]byte
	copy(keyBits[:], key)

	block, err := aes.NewCipher(keyBits[:])
	if err != nil {
		return nil, err
	}

	return &AESCipher{key: keyBits, block: block}, nil
}

func (c *AESCipher) Encrypt(payload string) ([]byte, error) {
	signed := ciphers.SignMessage(payload, c.key[:])

	data := []byte(signed)
	enc := make([]byte, aes.BlockSize+len(data))
	iv := enc[:aes.BlockSize]

	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(c.block, iv)
	stream.XORKeyStream(enc[aes.BlockSize:], data)

	return enc, nil
}

func (c *AESCipher) Decrypt(buf []byte) (string, error) {
	if len(buf) < aes.BlockSize {
		return "", errors.New("encrypted data is too small")
	}

	iv := buf[:aes.BlockSize]
	dec := buf[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(c.block, iv)
	stream.XORKeyStream(dec, dec)

	return ciphers.DecodeMessage(string(dec), c.key[:])
}
