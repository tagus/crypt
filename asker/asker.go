package asker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Asker struct {
	Reader *bufio.Reader
}

type Validation func(string) error

func DefaultAsker() *Asker {
	return &Asker{bufio.NewReader(os.Stdin)}
}

func (a *Asker) Ask(question string, validation Validation) (string, error) {
	if validation != nil {
		err := validation(question)
		if err != nil {
			return "", err
		}
	}

	fmt.Printf("%s ", question)
	ans, err := a.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return sanitize(ans), nil
}

func sanitize(str string) string {
	return strings.TrimSpace(str)
}
