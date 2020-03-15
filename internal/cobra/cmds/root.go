package cmds

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/store"
)

var (
	// Store is the current crypt store
	Store *store.CryptStore
	// Deving signals that current session is for development
	Deving bool
	// cryptfile refers to the path of the encrypted cryptfile
	cryptfile string
)

var rootCmd = &cobra.Command{
	Use:   "crypt",
	Short: "A secure credential store",
	Long: `Crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts.

Crypt uses a "cryptfile" to store any credentials securely. This file is
encrypted such that it cannot be read as plain text. There are a variety
of mechanisms to specify the crypt file, specified here in decreasing priority.

	1. cryptfile flag
	2. CRYPTFILE env variable
	3. ~/.crytpfile

===================================================================

Development mode offers an alternate path for a sample crypt file.
It does not prompt for a password. This is meant solely
for sandboxing. DO NOT STORE ANY CREDENTIALS HERE.`,
	SilenceUsage: true,
	// SilenceErrors: true,
}

// Execute executes the root cobra command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initCrypt)
	rootCmd.PersistentFlags().StringVarP(&cryptfile, "cryptfile", "c", "", "the cryptfile location")
	rootCmd.PersistentFlags().BoolVarP(&Deving, "dev", "d", false, "toggle development mode")
}

func printAndExit(err error) {
	if err != nil {
		// color.RedString(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

// getCryptfile determines the path of the cryptfile to be used, the cryptfile
// flag takes priority, falling back to a CRYPTFILE env var, and finally defaulting
// to a .cryptfile in the current user's home directory
func getCryptfile() (string, error) {
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

func initCrypt() {
	var (
		path, pwd string
		err       error
	)

	if Deving {
		pwd = "fakefakefake" // NOTE: development pwd, completely meaningless
		path = ".dev_cryptfile"
	} else {
		path, err = getCryptfile()
		printAndExit(err)
		asker := asker.DefaultAsker()
		secret, err := asker.AskSecret(color.YellowString("Password:"), false)
		printAndExit(err)
		pwd = secret
	}

	store, err := store.InitDefaultStore(path, pwd)
	printAndExit(err)

	Store = store
	color.Green("%s\n", "Crypt initialized successfully")
}

func serviceIsValid(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires exactly one arg")
	}
	if Store.Crypt.IsValid(args[0]) {
		return nil
	}
	suggestions := Store.Crypt.GetSuggestions(args[0])
	if len(suggestions) > 0 {
		fmt.Println("Invalid Service. Did you mean these instead?")
		for _, s := range suggestions {
			fmt.Printf("\t+ %s\n", s)
		}
	}
	return fmt.Errorf("invalid service specified")
}

func serviceIsNew(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires exactly one arg")
	}
	if !Store.Crypt.IsValid(args[0]) {
		return nil
	}
	return fmt.Errorf("service already exists")
}
