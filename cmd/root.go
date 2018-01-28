package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/backend"
)

var (
	Debugging   bool
	Credentials []backend.Credential
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
	home, _ := homedir.Dir()
	filePath := fmt.Sprintf("%s/.crypt", home)

	fmt.Println(filePath)
}
