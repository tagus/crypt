package ciphers

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputingHash(t *testing.T) {
	hash := ComputeHash("secret100")
	expected := "6c091667070fa0d70b8bab92755dec7af8d9adce498d08e18c6446e0d71d6cd9"
	assert.Equal(t, expected, hash)
}

func TestSigningMessage(t *testing.T) {
	signedMessage, err := SignMessage("Hello World", []byte("secret200"))
	require.NoError(t, err)
	require.Greater(t, len(signedMessage), 0)
}

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{
			name: "simple string",
			msg:  "Hello World",
		},
		{
			name: "json with delimiter",
			msg:  `{"test_name":"test_value"}`,
		},
	}

	key := []byte("secret200")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			signedMessage, err := SignMessage(test.msg, key)
			require.NoError(t, err)
			decodedMsg, err := DecodeMessage(signedMessage, key)
			assert.NoError(t, err)
			assert.Equal(t, test.msg, decodedMsg)
		})
	}
}

func TestDecodeMessageWhenMessageHasChanged(t *testing.T) {
	msg := "Hello World"
	key := []byte("secret200")
	signedMessage, err := SignMessage(msg, key)
	require.NoError(t, err)

	signedMessage[0] = 'A'
	_, err = DecodeMessage(signedMessage, key)
	assert.Error(t, err)
}
