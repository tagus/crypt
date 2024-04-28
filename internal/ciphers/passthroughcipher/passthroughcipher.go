package passthroughcipher

// PassThroughCipher provides no-op encryption and decryption and
// only be used as a placeholder or for testing
type PassThroughCipher struct{}

func New() *PassThroughCipher {
	return &PassThroughCipher{}
}

func (c *PassThroughCipher) Encrypt(val string) ([]byte, error) {
	return []byte(val), nil
}

func (c *PassThroughCipher) Decrypt(buf []byte) (string, error) {
	return string(buf), nil
}
