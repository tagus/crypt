package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"github.com/sugatpoudel/crypt/creds"
)

// An AES implementation for crypto with key size of 256 bits
type AesCrypto struct {
	key   [32]byte
	block cipher.Block
}

func InitAesCrypto(pwd string) (*AesCrypto, error) {
	key := computeHash(pwd)
	var keyBits [32]byte
	copy(keyBits[:], []byte(key))

	block, err := aes.NewCipher(keyBits[:])
	if err != nil {
		return nil, err
	}

	return &AesCrypto{keyBits, block}, nil
}

func (c *AesCrypto) Encrypt(crypt *creds.Crypt) ([]byte, error) {
	data, err := crypt.GetJson()
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

func (c *AesCrypto) Decrypt(enc []byte) (*creds.Crypt, error) {
	if len(enc) < aes.BlockSize {
		return nil, errors.New("Encrypted data is too small")
	}

	iv := enc[:aes.BlockSize]
	dec := enc[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(c.block, iv)
	stream.XORKeyStream(dec, dec)

	crypt, err := creds.FromJson(dec)
	if err != nil {
		return nil, err
	}

	return crypt, nil
}
