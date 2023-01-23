package secure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

// Computes an SHA 256 bit given the key string.
func computeHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashStr := hex.EncodeToString(hash.Sum(nil))
	return hashStr
}

// SignMessage signs the given message using hmac and appends it to
// the given message
func SignMessage(msg string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(msg))
	token := mac.Sum(nil)
	return msg + "-" + string(token)
}

// DecodeMessage decodes the given message by separating out the token
// and comparing it with the derived token, returning an error if the
// message was invalid
func DecodeMessage(msg string, key []byte) (string, error) {
	args := strings.Split(msg, "-")
	if len(args) != 2 {
		return "", errors.New("invalid message")
	}

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(args[0]))
	token := mac.Sum(nil)

	if hmac.Equal(token, []byte(args[1])) {
		return args[0], nil
	}
	return "", errors.New("invalid token")
}
