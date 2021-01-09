package cobracli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

// serviceIsValid determines if the given service is part of the current st
func serviceIsValid(cmd *cobra.Command, args []string) error {
	// TODO: support multi-word services
	if len(args) != 1 {
		return xerrors.New("requires exactly one arg")
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
		fmt.Println("invalid Service. Did you mean these instead?")
		for _, s := range suggestions {
			fmt.Printf("\t+ %s\n", s)
		}
	}
	return xerrors.Errorf("invalid service specified: %s", args[0])
}

// serviceIsNew ensures that the given service does not already exist
func serviceIsNew(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return xerrors.New("requires exactly one arg")
	}
	st, err := getStore()
	if err != nil {
		return err
	}
	if !st.Crypt.IsValid(args[0]) {
		return nil
	}
	return xerrors.Errorf("service already exists")
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
