package aescipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"github.com/tagus/crypt/internal/ciphers"
	"golang.org/x/crypto/bcrypt"
)

// AESCipher uses the AES encryption algorithm to encrypt and decrypt the provided data
// while ensuring that the data is signed with a hash of the provided password in order
// to ensure data integrity and provide password validation
type AESCipher struct {
	key       [32]byte
	block     cipher.Block
	signature []byte
}

func New(pwd string, hashedPwd, signature []byte) (*AESCipher, error) {
	err := bcrypt.CompareHashAndPassword(hashedPwd, []byte(pwd))
	if err != nil {
		return nil, ciphers.ErrInvalidPassword
	}

	key := ciphers.ComputeHash(pwd)
	var keyBits [32]byte
	copy(keyBits[:], key)

	block, err := aes.NewCipher(keyBits[:])
	if err != nil {
		return nil, err
	}

	return &AESCipher{
		key:       keyBits,
		block:     block,
		signature: signature,
	}, nil
}

func (c *AESCipher) Encrypt(payload string) ([]byte, error) {
	data, err := ciphers.SignMessage(payload, c.signature)
	if err != nil {
		return nil, err
	}

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

	return ciphers.DecodeMessage(dec, c.signature)
}
