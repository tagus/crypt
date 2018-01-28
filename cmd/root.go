package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/backend"
	input "github.com/tcnksm/go-input"
)

var (
	Debugging bool
	CryptFile backend.CryptFile
)

var rootCmd = &cobra.Command{
	Use:   "Crypt",
	Short: "A secure credential store",
	Long: `Crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initCrypt)
	rootCmd.PersistentFlags().BoolVarP(&Debugging, "debug", "d", false, "toggle debugging mode")
}

type Credential struct {
	Service  string
	Username string
	Password string
}

type Crypt struct {
	Credentials []Credential
}

func initCrypt() {
	homePath, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	cryptPath := filepath.Join(homePath, ".cryptfile")

	_, err = backend.ReadCrypt(cryptPath)
	if err != nil {
		newCryptFlow(ui, cryptPath)
	} else {

	}
}

func newCryptFlow(ui *input.UI, cryptPath string) {
	query := "No Crypt file found. Would you like to make one?"
	makeFile, err := ui.Ask(query, &input.Options{
		Default: "n",
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if makeFile == "yes" || makeFile == "y" {
		err := backend.MakeNewCrypt(cryptPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		color.Red("Did not create cryptfile")
		os.Exit(0)
	}
}
