package secure

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tagus/crypt/internal/creds"
)

var (
	ebay = creds.Credential{
		Service:     "eBay",
		Description: "electronic auction bay",
		Username:    "beanie_babies123",
		Password:    "mars321",
	}
	crypt = &creds.Crypt{
		Credentials: map[string]creds.Credential{
			"ebay": ebay,
		},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
)

func TestAesCryptoInit(t *testing.T) {
	_, err := InitAesCrypto("secret100")
	assert.Nil(t, err, "Aes Crypto was not initialized properly")
}

func TestEncryption(t *testing.T) {
	aes, _ := InitAesCrypto("secret100")
	_, err := aes.Encrypt(crypt)
	assert.Nil(t, err, "Encryption failed")
}

func TestDecryption(t *testing.T) {
	aes, _ := InitAesCrypto("secret100")
	enc, _ := aes.Encrypt(crypt)
	dec, err := aes.Decrypt(enc)
	assert.Nil(t, err, "Decryption failed")
	assert.Equal(t, crypt, dec, "Decryption was incorrect")
}
