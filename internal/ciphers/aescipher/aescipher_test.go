package aescipher

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tagus/crypt/internal/ciphers"
)

func TestAESCipher_Encrypt(t *testing.T) {
	signature := []byte("signature")
	hashedPwd, err := ciphers.ComputeHashPwd("secret100")
	require.NoError(t, err)

	ci, err := New("secret100", hashedPwd, signature)
	require.NoError(t, err)

	buf, err := ci.Encrypt("hello world")
	require.NoError(t, err)
	require.NotEmpty(t, buf)
}

func TestAESCipher_Decrypt(t *testing.T) {
	signature := []byte("signature")
	hashedPwd, err := ciphers.ComputeHashPwd("secret100")
	require.NoError(t, err)

	ci, err := New("secret100", hashedPwd, signature)
	require.NoError(t, err)

	buf, err := ci.Encrypt("hello world")
	require.NoError(t, err)
	require.NotEmpty(t, buf)

	dec, err := ci.Decrypt(buf)
	require.NoError(t, err)
	require.Equal(t, "hello world", dec)
}

func TestAESCipher_DecryptWithInvalidPassword(t *testing.T) {
	signature := []byte("signature")
	hashedPwd, err := ciphers.ComputeHashPwd("secret100")
	require.NoError(t, err)

	ci, err := New("secret100", hashedPwd, signature)
	require.NoError(t, err)

	buf, err := ci.Encrypt("hello world")
	require.NoError(t, err)
	require.NotEmpty(t, buf)

	_, err = New("secret200", hashedPwd, signature)
	require.ErrorIs(t, err, ErrInvalidPassword)
}

func TestAESCipher_DecryptWithInvalidSignature(t *testing.T) {
	signature := []byte("signature")
	hashedPwd, err := ciphers.ComputeHashPwd("secret100")
	require.NoError(t, err)

	ci, err := New("secret100", hashedPwd, signature)
	require.NoError(t, err)

	buf, err := ci.Encrypt("hello world")
	require.NoError(t, err)
	require.NotEmpty(t, buf)

	ci, err = New("secret100", hashedPwd, []byte("invalid-signature"))
	require.NoError(t, err)
	_, err = ci.Decrypt(buf)
	require.ErrorIs(t, err, ciphers.ErrInvalidSignature)
}
