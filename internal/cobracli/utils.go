package cobracli

import (
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/finder"
	"golang.org/x/xerrors"
)

// parseService parses the args to a single service credential and storing it globally
// so specific commands will not need to parse it again. This gives us a consolidated
// place to choose the right service.
func parseService(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return xerrors.Errorf("no service provided")
	}

	st, err := getStore()
	if err != nil {
		return err
	}

	fd := finder.New(st.Crypt)
	matches := fd.Filter(strings.Join(args, " "))
	if len(matches) == 0 {
		return xerrors.Errorf("invalid service provided: %v", args)
	}
	if len(matches) == 1 {
		return setService(matches[0])
	}

	possibles := make([]string, len(matches))
	for i := range matches {
		possibles[i] = matches[i].Service
	}

	asker := asker.DefaultAsker()
	n, err := asker.AskSelect("which service did you mean?", possibles)
	if err != nil {
		return err
	}
	return setService(matches[n])
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
	fd := finder.New(st.Crypt)
	cred := fd.Find(args[0])
	if cred == nil {
		return nil
	}
	return xerrors.Errorf("service already exists")
}

// saveStore persists any current changes to the st to the specified cryptfile
func saveStore() error {
	color.Green("saving crypt")
	st, err := getStore()
	if err != nil {
		return err
	}
	return st.Save()
}
