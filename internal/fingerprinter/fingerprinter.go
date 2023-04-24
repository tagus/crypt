package fingerprinter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/tagus/crypt/internal/crypt"
)

func Credential(c *crypt.Credential) (string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return hash(payload)
}

func Crypt(c *crypt.Crypt) (string, error) {
	payload, err := c.GetJSON()
	if err != nil {
		return "", err
	}
	return hash(payload)
}

func hash(payload []byte) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write(payload); err != nil {
		return "", err
	}
	hashStr := hex.EncodeToString(hash.Sum(nil))
	return hashStr, nil
}
