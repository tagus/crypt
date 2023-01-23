package cobracli

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/crypt"
	"github.com/tagus/crypt/internal/store"
)

var (
	st        *store.CryptStore
	svc       *crypt.Credential
	cryptfile string
)

var rootCmd = &cobra.Command{
	Use:   "crypt",
	Short: "a secure credential store",
	Long: `crypt is CLI application to securely store your credentials
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
	Version: crypt.Version,
}

// Execute executes the root cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if err == asker.ErrInterrupt {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cryptfile, "cryptfile", "c", "", "the cryptfile location")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(pwdCmd)
	rootCmd.AddCommand(showCmd)
}

// resolveCryptfilePath determines the path of the cryptfile to be used, the cryptfile
// flag takes priority, falling back to a CRYPTFILE env var, .cryptfile is the current working
// directory and finally defaulting to a .cryptfile in the current user's home directory
func resolveCryptfilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	hd, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	paths := []string{
		cryptfile,
		os.Getenv("CRYPTFILE"),
		filepath.Join(wd, ".cryptfile"),
		filepath.Join(hd, ".cryptfile"),
	}

	for _, path := range paths {
		if path == "" {
			continue
		}
		_, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		return path, nil
	}

	return "", errors.New("no valid cryptfile found")
}

func initStore() (*store.CryptStore, error) {
	path, err := resolveCryptfilePath()
	if err != nil {
		return nil, err
	}

	asker := asker.DefaultAsker()
	pwd, err := asker.AskSecret(color.YellowString("pwd"), false)
	if err != nil {
		return nil, err
	}

	return store.Decrypt(path, pwd)
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

func setService(service *crypt.Credential) error {
	if svc != nil {
		return errors.New("service already set")
	}
	svc = service
	return nil
}

func getService() (*crypt.Credential, error) {
	if svc == nil {
		return nil, errors.New("no service set")
	}
	return svc, nil
}
