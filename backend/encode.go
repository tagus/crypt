package backend

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/fatih/color"
)

// Encode encodes the given list of Crendtials into
// a byte array
func Encode(keystring string, cryptFile CryptFile) ([]byte, error) {
	json, err := json.Marshal(cryptFile)
	if err != nil {
		return nil, err
	}

	encryptedJson, err := Encrypt(keystring, json)
	if err != nil {
		return nil, err
	}

	return encryptedJson, nil
}

func Encrypt(keystring string, cryptJson []byte) ([]byte, error) {
	fmt.Println(string(cryptJson))
	key := []byte(keystring)

	// Creating the AES cypher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Empty array of cryptJson length with padding
	cipherText := make([]byte, aes.BlockSize+len(cryptJson))

	// IV included in the beginning
	iv := cipherText[:aes.BlockSize]

	// write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Get the encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(cipherText[aes.BlockSize:], cryptJson)

	return cipherText, nil
}

func SaveCrypt(cypherBytes []byte, cryptPath string) error {
	err := ioutil.WriteFile(cryptPath, cypherBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func MakeNewCrypt(keystring, cryptPath string) error {
	cryptFile := CryptFile{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cypherBytes, err := Encode(keystring, cryptFile)

	if err != nil {
		return err
	}

	err = SaveCrypt(cypherBytes, cryptPath)
	if err != nil {
		return err
	}

	fmt.Printf("cryptfile created at '%s'\n", color.YellowString("%s", cryptPath))
	return nil
}
