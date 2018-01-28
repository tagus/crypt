package cmd

import (
	"errors"
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
	CryptFile *backend.CryptFile
)

var rootCmd = &cobra.Command{
	Use:   "Crypt",
	Short: "A secure credential store",
	Long: `Crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts`,
}

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

	ui := input.DefaultUI()
	cryptPath := filepath.Join(homePath, ".cryptfile")

	fileBytes, err := backend.ReadCrypt(cryptPath)
	if err != nil {
		fmt.Println(err)
		newCryptFlow(ui, cryptPath)
	} else {
		keystring := readKey(ui)
		CryptFile, err = backend.Decode(keystring, fileBytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
		keyString := readKey(ui)
		err = backend.MakeNewCrypt(keyString, cryptPath)
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

func readKey(ui *input.UI) string {
	query := "What is your key?"
	keyString, err := ui.Ask(query, &input.Options{
		Required: true,
		Loop:     true,
		Mask:     false,
		ValidateFunc: func(s string) error {
			if len(s) < 16 {
				msg := color.RedString("%s", "key is too small. must be at least 16 bytes")
				return errors.New(msg)
			}

			return nil
		},
	})

	if err != nil {
		fmt.Println(keyString)
		os.Exit(1)
	}

	return keyString
}
