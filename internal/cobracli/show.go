package cobracli

import (
	"errors"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/crypt"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [service]",
	Short: "show information about a service",
	Long: `prints information about the specified service.
Provided service needs to already exist in crypt.`,
	Args: parseService,
	RunE: show,
}

func show(cmd *cobra.Command, args []string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	err = clipboard.WriteAll(svc.Password)
	if err != nil {
		return errors.New("failed to copy pwd")
	}

	crypt.PrintCredential(svc)
	return nil
}
