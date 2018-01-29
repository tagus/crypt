package crypto

import "github.com/sugatpoudel/crypt/creds"

// A simple interface for any struct that can encrypt
// and decrypt crypt data.
type Crypto interface {
	Encrypt(crypt *creds.Crypt) ([]byte, error)
	Decrypt(cipher []byte) (*creds.Crypt, error)
}
