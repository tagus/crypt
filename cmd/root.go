package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/store"
)

var (
	Store  *store.CryptStore
	Deving bool
)

var rootCmd = &cobra.Command{
	Use:   "crypt",
	Short: "A secure credential store",
	Long: `Crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts.

Crypt assumes the existence of a '.cryptfile' in the home directory
and tries to decrypt it upon initialization. If such file does not
exists, one will be created.

Development mode offers an alternate path for a sample crypt file.
It does not prompt for a password. This is meant solely
for sandboxing. DO NOT STORE ANY CREDENTIALS HERE.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initCrypt)
	rootCmd.PersistentFlags().BoolVarP(&Deving, "dev", "d", false, "toggle development mode")
}

func printAndExit(err error) {
	if err != nil {
		// color.RedString(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCrypt() {
	var pwd string
	var filename string
	if Deving {
		pwd = "fakefakefake" // NOTE: development pwd, completely meaningless
		filename = ".dev_cryptfile"
	} else {
		fmt.Printf("%s ", color.YellowString("Password:"))
		pwdB, err := gopass.GetPasswd()
		printAndExit(err)

		pwd = string(pwdB)
		filename = ".cryptfile"
	}

	home, err := homedir.Dir()
	printAndExit(err)

	path := filepath.Join(home, filename)
	store, err := store.InitDefaultStore(path, string(pwd))
	printAndExit(err)

	Store = store
	color.Green("%s\n\n", "Crypt initialized successfully")
}

func serviceIsValid(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires exactly one arg")
	}
	if Store.Crypt.IsValid(args[0]) {
		return nil
	}
	return errors.New("invalid service specified")
}
