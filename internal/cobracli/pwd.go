package cobracli

import (
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/finder"
)

// pwdCmd represents the pwd command
var pwdCmd = &cobra.Command{
	Use:   "pwd [service]",
	Short: "get the password for a service",
	Long:  `copies the password for a service to the clipboard.`,
	Args:  serviceIsValid,
	RunE:  getPwd,
}

func getPwd(cmd *cobra.Command, args []string) error {
	service := args[0]

	st, err := getStore()
	if err != nil {
		return err
	}
	fd := finder.New(st.Crypt)
	cred := fd.Find(service)
	pwd := cred.Password

	err = clipboard.WriteAll(pwd)
	if err != nil {
		return err
	}

	color.Green("pwd copied to clipboard")
	return nil
}
