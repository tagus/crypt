package cobracli

import (
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
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
		return xerrors.Errorf("service was not found: %s", service)
	}

	err = clipboard.WriteAll(cred.Password)
	if err != nil {
		return xerrors.Errorf("failed to copy pwd")
	}

	cred.PrintCredential()
	return nil
}
