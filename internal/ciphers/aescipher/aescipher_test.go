package aescipher

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAESCipher_Encrypt(t *testing.T) {
	ci, err := New("secret100")
	require.NoError(t, err)

	buf, err := ci.Encrypt("hello world")
	require.NoError(t, err)
	require.NotEmpty(t, buf)
}

func TestAESCipher_Decrypt(t *testing.T) {
	ci, err := New("secret100")
	require.NoError(t, err)

	buf, err := ci.Encrypt("hello world")
	require.NoError(t, err)
	require.NotEmpty(t, buf)

	dec, err := ci.Decrypt(buf)
	require.NoError(t, err)
	require.Equal(t, "hello world", dec)
}

func TestAESCipher_DecryptWithInvalidPassword(t *testing.T) {
	ci, err := New("secret100")
	require.NoError(t, err)

	buf, err := ci.Encrypt("hello world")
	require.NoError(t, err)
	require.NotEmpty(t, buf)

	ciIncorrect, err := New("secret200")
	require.NoError(t, err)
	_, err = ciIncorrect.Decrypt(buf)
	require.Error(t, err)
}
