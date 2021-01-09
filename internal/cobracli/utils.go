package cobracli

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// serviceIsValid determines if the given service is part of the current st
func serviceIsValid(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires exactly one arg")
	}
	st, err := getStore()
	if err != nil {
		return err
	}
	if st.Crypt.IsValid(args[0]) {
		return nil
	}
	suggestions := st.Crypt.GetSuggestions(args[0])
	if len(suggestions) > 0 {
		// TODO: use selector to select from suggestions
		fmt.Println("Invalid Service. Did you mean these instead?")
		for _, s := range suggestions {
			fmt.Printf("\t+ %s\n", s)
		}
	}
	return fmt.Errorf("invalid service specified")
}

// serviceIsNew ensures that the given service does not already exist
func serviceIsNew(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires exactly one arg")
	}
	st, err := getStore()
	if err != nil {
		return err
	}
	if !st.Crypt.IsValid(args[0]) {
		return nil
	}
	return fmt.Errorf("service already exists")
}

// saveStore persists any current changes to the st to the specified cryptfile
func saveStore() error {
	color.Green("\nSaving crypt")
	st, err := getStore()
	if err != nil {
		return err
	}
	return st.Save()
}