package ciphers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputingHash(t *testing.T) {
	hash := ComputeHash("secret100")
	expected := "6c091667070fa0d70b8bab92755dec7af8d9adce498d08e18c6446e0d71d6cd9"
	assert.Equal(t, expected, hash, "Unequal hash")
}

func TestSigningMessage(t *testing.T) {
	signedMessage := SignMessage("Hello World", []byte("secret200"))
	assert.NotEqual(t, "", signedMessage)
}

func TestDecodeMessage(t *testing.T) {
	msg := "Hello World"
	key := []byte("secret200")
	signedMessage := SignMessage(msg, key)
	decodedMsg, err := DecodeMessage(signedMessage, key)
	assert.NoError(t, err, "An error occurred")
	assert.Equal(t, msg, decodedMsg, "Decoded message was not equal to original message")
}

func TestDecodeMessageWhenMessageHasChanged(t *testing.T) {
	msg := "Hello World"
	key := []byte("secret200")
	signedMessage := SignMessage(msg, key)
	changedMessage := strings.Replace(signedMessage, "Hello", "Goodbye", 1)
	_, err := DecodeMessage(changedMessage, key)
	assert.Error(t, err, "An error was expected")
}
