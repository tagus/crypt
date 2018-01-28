package backend

import (
	"io/ioutil"
)

// Decode decodes the given byte stream and converts it to
// a slice of Credentials
func Decode(fileBytes []byte) (CryptFile, error) {
	return CryptFile{}, nil
}

func ReadCrypt(filePath string) ([]byte, error) {
	_, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
