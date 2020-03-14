package cmds

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [service]",
	Short: "Show information about a service",
	Long: `Prints information about the specified service.
Provided service needs to already exist in crypt.`,
	Args: serviceIsValid,
	Run:  show,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func show(cmd *cobra.Command, args []string) {
	service := args[0]
	cred := Store.Crypt.FindCredential(service)
	if cred == nil {
		color.Red("service '%s' was not found", service)
		return
	}

	err := clipboard.WriteAll(cred.Password)
	if err != nil {
		fmt.Println(color.YellowString("failed to copy password"))
	}

	cred.PrintCredential()
}