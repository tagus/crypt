package cmd

import (
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// pwdCmd represents the pwd command
var pwdCmd = &cobra.Command{
	Use:   "pwd [service]",
	Short: "Get the password for a service",
	Long:  `Copies the password for a service to the clipboard.`,
	Args:  serviceIsValid,
	Run:   getPwd,
}

func init() {
	rootCmd.AddCommand(pwdCmd)
}

func getPwd(cmd *cobra.Command, args []string) {
	service := args[0]
	cred := Store.Crypt.FindCredential(service)
	pwd := cred.Password

	err := clipboard.WriteAll(pwd)
	printAndExit(err)

	color.Green("Password copied to clipboard")
}
