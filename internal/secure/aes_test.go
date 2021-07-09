package secure

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tagus/crypt/internal/crypt"
)

var (
	ebay = &crypt.Credential{
		Service:     "eBay",
		Description: "electronic auction bay",
		Username:    "beanie_babies123",
		Password:    "mars321",
	}
	cr = &crypt.Crypt{
		Credentials: map[string]*crypt.Credential{
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
	_, err := aes.Encrypt(cr)
	assert.Nil(t, err, "Encryption failed")
}

func TestDecryption(t *testing.T) {
	aes, _ := InitAesCrypto("secret100")
	enc, _ := aes.Encrypt(cr)
	dec, err := aes.Decrypt(enc)
	assert.Nil(t, err, "Decryption failed")
	assert.Equal(t, cr, dec, "Decryption was incorrect")
}
