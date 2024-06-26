package asker

import (
	"errors"
	"io"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/tagus/mango"
)

var (
	ErrInterrupt = errors.New("interrupted through signal")

	confirmKeywords = mango.NewSet("si", "y", "ack", "yup", "ok", "yes")
)

// Asker is a helper for retrieving user input from the given reader
type Asker struct {
	IsVimMode bool
	Mask      rune
	Stdin     io.ReadCloser
	Stdout    io.WriteCloser
}

// Validation defines a general validation processor for a given string
type Validation func(string) error

// DefaultAsker creates an asker using STDOUT and STDIN
func DefaultAsker() *Asker {
	return &Asker{
		Mask:      '*',
		IsVimMode: true,
	}
}

// validator is higher order function that creates a Validator from the
// given list of validation function.
func validator(validations ...Validation) promptui.ValidateFunc {
	return func(label string) error {
		for _, validation := range validations {
			err := validation(label)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Ask retrieves user input from the current Reader
func (a *Asker) Ask(question string, validations ...Validation) (string, error) {
	prompt := promptui.Prompt{
		Label:     question,
		Validate:  validator(validations...),
		Stdin:     a.Stdin,
		Stdout:    a.Stdout,
		IsVimMode: a.IsVimMode,
	}

	res, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return sanitize(res), nil
}

// AskConfirm retrieves confirmation and attempts to parse the various way
// that a use might supply confirmation e.g. yes, YES, y, si, yup etc.
func (a *Asker) AskConfirm(question string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     question,
		Stdin:     a.Stdin,
		Stdout:    a.Stdout,
		IsVimMode: a.IsVimMode,
	}

	res, err := prompt.Run()
	if err != nil {
		return false, err
	}

	res = strings.ToLower(res)
	return confirmKeywords.Has(res), nil
}

// AskSecret asks for user input without echoing input back to terminal.
// Note that this method is only supported through stdin
func (a *Asker) AskSecret(question string, confirm bool, validations ...Validation) (string, error) {
	ask := promptui.Prompt{
		Label:     question,
		Validate:  validator(validations...),
		Stdin:     a.Stdin,
		Stdout:    a.Stdout,
		IsVimMode: a.IsVimMode,
		Mask:      a.Mask,
	}

	res, err := ask.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrEOF) || errors.Is(err, promptui.ErrInterrupt) {
			return "", ErrInterrupt
		}
		return "", err
	}

	if confirm {
		_, err = a.AskSecret("confirm "+question, false, func(val string) error {
			if val != res {
				return errors.New("confirmation does not match")
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}

	return sanitize(res), nil
}

// AskSelect prompts the user from a list of items and returns the index of the selected item
func (a *Asker) AskSelect(question string, items any) (int, error) {
	ask := promptui.Select{
		Label:  question,
		Items:  items,
		Stdin:  a.Stdin,
		Stdout: a.Stdout,
	}
	index, _, err := ask.Run()
	if err != nil {
		return 0, err
	}
	return index, nil
}

// sanitize is a helper function to clean up the user input
func sanitize(str string) string {
	return strings.TrimSpace(str)
}
