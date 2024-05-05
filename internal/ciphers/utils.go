package ciphers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// ComputeHash computes a SHA 256 bit given the key string.
func ComputeHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashStr := hex.EncodeToString(hash.Sum(nil))
	return hashStr
}

// SignMessage signs the given message using hmac and appends it to
// the given message
func SignMessage(msg string, key []byte) ([]byte, error) {
	if key == nil || len(key) == 0 {
		return nil, errors.New("invalid key given")
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(msg))
	token := mac.Sum(nil)

	buf := []byte(msg)
	buf = append(buf, '&')
	buf = append(buf, token...)
	return buf, nil
}

// DecodeMessage decodes the given message by separating out the token
// and comparing it with the derived token, returning an error if the
// message was invalid
func DecodeMessage(buf []byte, key []byte) (string, error) {
	if len(buf) < sha256.Size {
		return "", errors.New("invalid signed message: message is too short to account for signature")
	}

	split := len(buf) - sha256.Size
	if buf[split-1] != '&' {
		return "", errors.New("invalid signed message: no delimiter found")
	}
	msg, signature := buf[:split-1], buf[split:]

	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	token := mac.Sum(nil)

	if hmac.Equal(token, signature) {
		return string(msg), nil
	}
	return "", errors.New("invalid token")
}
