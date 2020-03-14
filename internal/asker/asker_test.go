package asker

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// generate a temporary file to mock user input
func getTempFile() *os.File {
	temp, err := ioutil.TempFile("", "asker_test")
	if err != nil {
		panic(err)
	}

	return temp
}

func TestValidation(t *testing.T) {
	temp := getTempFile()
	asker := &Asker{bufio.NewReader(temp)}
	defer os.Remove(temp.Name())

	validation := func(input string) error {
		if len(input) < 10 {
			return errors.New("Input is too short")
		}
		return nil
	}

	_, err := asker.Ask("What is your number?\n", validation)
	assert.NotNil(t, err)

	temp.Close()
}

func TestAsker(t *testing.T) {
	tmp := getTempFile()
	asker := &Asker{bufio.NewReader(tmp)}
	defer os.Remove(tmp.Name())

	tmp.WriteString("Tagus Leduop\n")
	tmp.Seek(0, os.SEEK_SET)

	ans, err := asker.Ask("What's your name?\n", nil)

	assert.NotEmpty(t, ans)
	assert.Nil(t, err)

	tmp.Close()
}
