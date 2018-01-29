package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/files"
)

var Store *files.CryptStore

var rootCmd = &cobra.Command{
	Use:   "crypt",
	Short: "A secure credential store",
	Long: `Crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts.

Crypt assumes the existence of a '.cryptfile' in the home directory
and tries to decrypt it upon initialization. If such file does not
exists, one will be created.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initCrypt)
}

func printAndExit(err error) {
	if err != nil {
		// color.RedString(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCrypt() {
	fmt.Printf("%s ", color.YellowString("Password:"))
	pwd, err := gopass.GetPasswd()
	printAndExit(err)

	home, err := homedir.Dir()
	printAndExit(err)

	path := filepath.Join(home, ".cryptfile")
	store, err := files.InitDefaultStore(path, string(pwd))
	printAndExit(err)

	Store = store
	color.Green("%s\n\n", "Crypt initialized successfully")
}
