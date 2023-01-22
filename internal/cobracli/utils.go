package cobracli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
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

	fd, err := finder.New(st.Crypt)
	if err != nil {
		return err
	}

	q := strings.Join(args, " ")
	matches, err := fd.Filter(q)
	if err != nil {
		return err
	}
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

// backupStore duplicates the current store at the current location at the instant level
func backupCrypt(cmd *cobra.Command, args []string) error {
	st, err := getStore()
	if err != nil {
		return err
	}

	ts := time.Now().Unix()
	backupFile := fmt.Sprintf("%s-%d.cryptfile", st.Crypt.Id, ts)
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	backupDir := filepath.Join(home, ".cryptrc", "backups")
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return err
	}

	return st.SaveTo(filepath.Join(backupDir, backupFile))
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
	fd, err := finder.New(st.Crypt)
	if err != nil {
		return err
	}
	cred, err := fd.Find(args[0])
	if err != nil {
		return err
	}
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

func combineArgs(fns ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range fns {
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
