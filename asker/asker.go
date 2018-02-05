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

type Asker struct {
	Reader *bufio.Reader
}

type Validation func(string) error

func DefaultAsker() *Asker {
	return &Asker{bufio.NewReader(os.Stdin)}
}

func (a *Asker) Ask(question string, validation Validation) (string, error) {
	fmt.Printf("%s ", question)
	ans, err := a.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if validation != nil {
		err := validation(ans)
		if err != nil {
			return "", err
		}
	}

	return sanitize(ans), nil
}

// Asks for user input without echoing input back to terminal.
// Note that this method is only supported through stdin
func (a *Asker) AskSecret(question string, validation Validation) (string, error) {
	fmt.Printf("%s ", question)
	pwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	fmt.Printf("\n")

	fmt.Printf("Confirm %s ", question)
	conf, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	fmt.Printf("\n")

	if !bytes.Equal(pwd, conf) {
		return "", errors.New("Confirmation did not match.")
	}

	ans := string(pwd)
	if validation != nil {
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
