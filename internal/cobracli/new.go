package cobracli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/store"
	"golang.org/x/xerrors"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "create new cryptfile",
	Long:    "attempts to create a new cryptfile at the resolved path if one does not already exist.",
	RunE:    new,
	Aliases: []string{"init"},
}

func new(cmd *cobra.Command, args []string) error {
	path, err := resolveCryptfilePath()
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			asker := asker.DefaultAsker()
			pwd, err := asker.AskSecret("enter pwd", true)
			if err != nil {
				return err
			}
			_, err = store.InitDefaultStore(path, pwd)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return xerrors.New("cryptfile already exists")
}
