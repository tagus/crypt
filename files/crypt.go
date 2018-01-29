package files

import (
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/sugatpoudel/crypt/creds"
	"github.com/sugatpoudel/crypt/secure"
)

const perm = 0600

// Represents a crypt instance stored as a file
type CryptStore struct {
	path   string
	crypto secure.Crypto
	Crypt  *creds.Crypt
}

// Creates an empty crypt file in the given path.
func createDefaultCryptFile(path string) error {
	credMap := make(map[string]creds.Credential)
	now := time.Now().Unix()
	crypt := &creds.Crypt{credMap, now, now}

	json, err := crypt.GetJson()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, json, perm)
	if err != nil {
		return err
	}

	return nil
}

// Initializes a default crypt store using the AES crypto implementation.
// If the crypt file does not exist, one will be created in the provided path.
func InitDefaultStore(path, pwd string) (*CryptStore, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := createDefaultCryptFile(path)
		if err != nil {
			return nil, err
		}
	}

	crypto, err := secure.InitAesCrypto(pwd)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	crypt, err := crypto.Decrypt(data)
	if err != nil {
		return nil, errors.New("Password was invalid. Decryption failed.")
	}

	store := &CryptStore{path, crypto, crypt}
	return store, nil
}

// Encrypts the current Crypt and saves it to the path field.
func (s *CryptStore) Save() error {
	data, err := s.crypto.Encrypt(s.Crypt)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.path, data, perm)
	if err != nil {
		return err
	}

	return nil
}

// Recreates the Crypto instance with the new password.
func (s *CryptStore) ChangePwd(pwd string) error {
	crypto, err := secure.InitAesCrypto(pwd)
	if err != nil {
		return err
	}

	s.crypto = crypto
	return nil
}
