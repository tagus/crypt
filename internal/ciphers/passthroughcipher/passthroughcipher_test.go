package passthroughcipher

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPassThroughCipher_Encrypt(t *testing.T) {
	ci := New()
	encrypted, err := ci.Encrypt("test")
	require.NoError(t, err)
	require.Equal(t, []byte("test"), encrypted)
}

func TestPassThroughCipher_Decrypt(t *testing.T) {
	ci := New()
	decrypted, err := ci.Decrypt([]byte("test"))
	require.NoError(t, err)
	require.Equal(t, "test", decrypted)
}
