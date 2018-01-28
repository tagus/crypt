package backend

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Decode decodes the given byte stream and converts it to
// a slice of Credentials
func Decode(keystring string, cipherCrypt []byte) (*CryptFile, error) {
	cryptFileBytes, err := Decrypt(keystring, cipherCrypt)
	if err != nil {
		return nil, err
	}

	var cryptFile CryptFile
	err = json.Unmarshal(cryptFileBytes, &cryptFile)
	if err != nil {
		return nil, err
	}

	return &cryptFile, nil
}

func Decrypt(keystring string, cipherCrypt []byte) ([]byte, error) {
	key := []byte(keystring)

	// Create the AES cypher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// sanity check for description
	if len(cipherCrypt) < aes.BlockSize {
		return nil, errors.New("Cryptfile is too small")
	}

	// get the 16 byte IV
	iv := cipherCrypt[:aes.BlockSize]

	// remove the IV from the cipherCrypt
	cipherCrypt = cipherCrypt[aes.BlockSize:]

	// get the decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// decrypt bytes from cipherCrypt
	stream.XORKeyStream(cipherCrypt, cipherCrypt)

	return cipherCrypt, nil
}

func ReadCrypt(filePath string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
