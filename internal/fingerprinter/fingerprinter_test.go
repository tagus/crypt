package fingerprinter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tagus/crypt/internal/crypt"
)

var (
	ebay1 = &crypt.Credential{
		Id:          "acct-1",
		Service:     "eBay",
		Description: "electronic auction bay",
		Username:    "beanie_babies123",
		Password:    "mars321",
	}
	ebay2 = &crypt.Credential{
		Id:          "acct-1",
		Service:     "eBay",
		Description: "electronic auction bay",
		Username:    "beanie_babies123",
		Password:    "mars321",
	}
	amazon = &crypt.Credential{
		Id:          "acct-2",
		Service:     "Amazon Web Services",
		Description: "america's ali baba",
		Email:       "jeff.bezos@amazon.com",
		Password:    "123jupiter",
	}
)

func TestFingerprinter_Credential(t *testing.T) {
	ebay1Hash, err := Credential(ebay1)
	assert.NoError(t, err)

	ebay2Hash, err := Credential(ebay2)
	assert.NoError(t, err)

	amazonHash, err := Credential(amazon)
	assert.NoError(t, err)

	assert.Equal(t, ebay1Hash, ebay2Hash)
	assert.NotEqual(t, amazonHash, ebay2Hash)
}

func TestFingerprinter_Crypt(t *testing.T) {
	c := &crypt.Crypt{
		Id:      "crypt-1",
		Version: "v1",
		Credentials: crypt.Credentials{
			"acct-1": ebay1,
			"acct-2": amazon,
		},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	fingerprint, err := Crypt(c)
	assert.NoError(t, err)
	assert.NotEmpty(t, fingerprint)

	c.Fingerprint = fingerprint

	fingerprint2, err := Crypt(c)
	assert.NoError(t, err)
	assert.NotEqual(t, fingerprint2, fingerprint)
}
