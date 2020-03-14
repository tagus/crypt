package asker

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// Asker is a helper for retrieving user input from the given reader
type Asker struct {
	Reader *bufio.Reader
}

// Validation defines a general validation processor for a given string
type Validation func(string) error

// DefaultAsker constructs a asker from stdin
func DefaultAsker() *Asker {
	return &Asker{bufio.NewReader(os.Stdin)}
}

// Ask retrieves user input from the current Reader
func (a *Asker) Ask(question string, validations ...Validation) (string, error) {
	fmt.Printf("%s ", question)
	ans, err := a.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	for _, validation := range validations {
		err := validation(ans)
		if err != nil {
			return "", err
		}
	}

	return sanitize(ans), nil
}

// AskSecret asks for user input without echoing input back to terminal.
// Note that this method is only supported through stdin
func (a *Asker) AskSecret(question string, confirm bool, validations ...Validation) (string, error) {
	fmt.Printf("%s ", question)
	pwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	fmt.Printf("\n")

	if confirm {
		fmt.Printf("Confirm %s ", question)
		conf, err := terminal.ReadPassword(0)
		if err != nil {
			return "", err
		}
		fmt.Printf("\n")

		if !bytes.Equal(pwd, conf) {
			return "", errors.New("confirmation did not match")
		}
	}

	ans := string(pwd)
	for _, validation := range validations {
		err := validation(ans)
		if err != nil {
			return "", err
		}
	}

	return sanitize(ans), nil
}

func sanitize(str string) string {
	return strings.TrimSpace(str)
}
