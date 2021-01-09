package cobracli

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [service]",
	Short: "show information about a service",
	Long: `prints information about the specified service.
Provided service needs to already exist in crypt.`,
	Args: serviceIsValid,
	RunE: show,
}

func show(cmd *cobra.Command, args []string) error {
	service := args[0]

	st, err := getStore()
	if err != nil {
		return err
	}

	cred := st.FindCredential(service)
	if cred == nil {
		return fmt.Errorf("service was not found: %s", service)
	}

	err = clipboard.WriteAll(cred.Password)
	if err != nil {
		return fmt.Errorf("failed to copy pwd")
	}

	cred.PrintCredential()
	return nil
}
