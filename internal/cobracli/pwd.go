package cobracli

import (
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// pwdCmd represents the pwd command
var pwdCmd = &cobra.Command{
	Use:   "pwd [service]",
	Short: "get the password for a service",
	Long:  `copies the password for a service to the clipboard.`,
	Args:  parseService,
	RunE:  getPwd,
}

func getPwd(cmd *cobra.Command, args []string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	pwd := svc.Password
	err = clipboard.WriteAll(pwd)
	if err != nil {
		return err
	}

	color.Green("pwd copied to clipboard")
	return nil
}
