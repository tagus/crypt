package legacycrypt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ebay = &Credential{
		Id:          "acct-1",
		Service:     "eBay",
		Description: "electronic auction bay",
		Username:    "beanie_babies123",
		Password:    "mars321",
	}
	amazon = &Credential{
		Id:          "acct-2",
		Service:     "Amazon Web Services",
		Description: "america's ali baba",
		Email:       "jeff.bezos@amazon.com",
		Password:    "123jupiter",
	}

	crypt = &Crypt{
		Credentials: Credentials{
			"acct-1": ebay,
		},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
)

func TestCrypt_SetCredential(t *testing.T) {
	crypt.SetCredential(amazon)
	cred := crypt.GetCredential("acct-2")
	assert.NotNil(t, cred, "Did not set credential properly")
}

func TestCrypt_SetCredentialWithoutId(t *testing.T) {
	cred, err := crypt.SetCredential(&Credential{
		Service:     "foo service",
		Description: "foo description",
		Email:       "foo@bar.baz",
		Password:    "fizzbuzz",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, cred.Id)
}

func TestCrypt_RemoveCredential(t *testing.T) {
	crypt.RemoveCredential("acct-1")
	cred := crypt.GetCredential("acct-1")
	assert.Nil(t, cred, "Did not remove credential properly")
}

func TestCrypt_ToAndFromJSON(t *testing.T) {
	buf, err := crypt.GetJSON()
	assert.Nil(t, err, "Did not marshall properly")

	actual, err := FromJSON(buf)
	assert.Nil(t, err, "failed to unmarshal from json")
	assert.NotNil(t, actual, "unmarshalled crypt cannot be nil")
}
