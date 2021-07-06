package cobracli

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/creds"
	"github.com/tagus/crypt/internal/store"
	"golang.org/x/xerrors"
)

var seedFile string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "create new cryptfile",
	Long:    "attempts to create a new cryptfile at the resolved path if one does not already exist.",
	RunE:    new,
	Aliases: []string{"init"},
}

func init() {
	newCmd.Flags().StringVarP(&seedFile, "seed", "s", "", "a plaintext crypt file to seed from")
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

			crypt, err := buildCrypt()
			if err != nil {
				return err
			}

			_, err = store.InitDefaultStore(path, pwd, crypt)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return xerrors.New("cryptfile already exists")
}

func buildCrypt() (*creds.Crypt, error) {
	if seedFile == "" {
		credMap := make(map[string]*creds.Credential)
		now := time.Now().Unix()
		return &creds.Crypt{
			Credentials: credMap,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil
	}

	buf, err := ioutil.ReadFile(seedFile)
	if err != nil {
		return nil, err
	}
	return creds.FromJSON(buf)
}
