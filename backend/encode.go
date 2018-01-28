package backend

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// Encode encodes the given list of Crendtials into
// a byte array
func Encode(cryptFile CryptFile) ([]byte, error) {
	return nil, nil
}

func SaveCrypt(cypherBytes []byte) error {
	return nil
}

func MakeNewCrypt(cryptPath string) error {
	cryptFile := CryptFile{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cypherBytes, err := Encode(cryptFile)

	if err != nil {
		return err
	}

	err = SaveCrypt(cypherBytes)
	if err != nil {
		return err
	}

	fmt.Printf("cryptfile created at '%s'\n", color.YellowString("%s", cryptPath))
	return nil
}
