package cmds

import (
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/utils"
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
	cred := getStore().Crypt.FindCredential(service)
	pwd := cred.Password

	err := clipboard.WriteAll(pwd)
	utils.FatalIf(err)

	color.Green("Password copied to clipboard")
}
