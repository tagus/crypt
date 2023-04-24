package cobracli

import (
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/crypt"
	"github.com/tagus/crypt/internal/fingerprinter"
	"github.com/tagus/crypt/internal/store"
	"github.com/teris-io/shortid"
)

var seedFile string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new [path]",
	Short:   "create new cryptfile",
	Long:    "attempts to create a new cryptfile at the resolved path if one does not already exist.",
	RunE:    new,
	Aliases: []string{"init"},
	Args:    cobra.ExactArgs(1),
}

func init() {
	newCmd.Flags().StringVarP(&seedFile, "seed", "s", "", "a plaintext crypt file to seed from")
}

func new(cmd *cobra.Command, args []string) error {
	path := args[0]
	_, err := os.Stat(path)
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
	return errors.New("cryptfile already exists")
}

func buildCrypt() (*crypt.Crypt, error) {
	if seedFile == "" {
		id, err := shortid.Generate()
		if err != nil {
			return nil, err
		}

		credMap := make(crypt.Credentials)
		now := time.Now().Unix()
		cr := &crypt.Crypt{
			Id:          id,
			Credentials: credMap,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		fp, err := fingerprinter.Crypt(cr)
		if err != nil {
			return nil, err
		}
		cr.Fingerprint = fp
		return cr, nil
	}

	buf, err := ioutil.ReadFile(seedFile)
	if err != nil {
		return nil, err
	}
	return crypt.FromJSON(buf)
}
