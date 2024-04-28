package passthroughcipher

// PassthroughCipher provides no-op encryption and decryption and
// only be used as a placeholder or for testing
type PassthroughCipher struct{}

func New() *PassthroughCipher {
	return &PassthroughCipher{}
}

func (c *PassthroughCipher) Encrypt(val string) ([]byte, error) {
	return []byte(val), nil
}

func (c *PassthroughCipher) Decrypt(buf []byte) (string, error) {
	return string(buf), nil
}
