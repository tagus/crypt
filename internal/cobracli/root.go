package cobracli

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/env"
	"github.com/sugatpoudel/crypt/internal/store"
	"github.com/sugatpoudel/crypt/internal/utils"
)

var (
	st        *store.CryptStore
	deving    bool
	cryptfile string
)

var rootCmd = &cobra.Command{
	Use:   "crypt",
	Short: "A secure credential store",
	Long: `Crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts.

Although crypt supports multi word names for a credential, it might be cumbersome
to retrieve it in the future. Thus, it might be easier to stick to dash separated
one word names.

Crypt uses a "cryptfile" to store any credentials securely. This file is
encrypted such that it cannot be read as plain text. There are a variety
of mechanisms to specify the crypt file, specified here in decreasing priority.

	1. cryptfile flag
	2. CRYPTFILE env variable
	3. ~/.crytpfile`,
	SilenceUsage: true,
	// SilenceErrors: true,
	Version: "v0.1.1",
}

// Execute executes the root cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cryptfile, "cryptfile", "c", "", "the cryptfile location")
	rootCmd.PersistentFlags().BoolVarP(&deving, "dev", "d", false, "toggle development mode")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(pwdCmd)
	rootCmd.AddCommand(showCmd)
}

// getStore is a helper function to retrieve the current crypt store
// this method will initialize a crypt store if one does not already exist.
func getStore() (*store.CryptStore, error) {
	if st == nil {
		store, err := initStore()
		if err != nil {
			return nil, err
		}
		st = store
	}
	return st, nil
}

func initStore() (*store.CryptStore, error) {
	var (
		path, pwd string
		err       error
	)
	if deving {
		// TODO: see if this flag can be disabled in prod.
		// Note: this sets up a dev cryptfile where the crypt cmd was executed with a
		// meaningless password. This is meant to help quickly test features without
		// going through an auth wall. Do not store actual credentials here.
		env.SetDev(true)
		pwd = "fakefakefake"
		path = ".dev_cryptfile"
	} else {
		path, err = resolveCryptfilePath()
		utils.FatalIf(err)

		asker := asker.DefaultAsker()
		secret, err := asker.AskSecret(color.YellowString("Password"), false)
		utils.FatalIf(err)

		pwd = secret
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("cryptfile does not exist")
		}
		return nil, err
	}

	return store.InitDefaultStore(path, pwd)
}

// resolveCryptfilePath determines the path of the cryptfile to be used, the cryptfile
// flag takes priority, falling back to a CRYPTFILE env var, and finally defaulting
// to a .cryptfile in the current user's home directory
func resolveCryptfilePath() (string, error) {
	if cryptfile != "" {
		return cryptfile, nil
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	path, ok := os.LookupEnv("CRYPTFILE")
	if ok {
		return path, nil
	}
	path = filepath.Join(home, ".cryptfile")
	return path, nil
}
