package cobracli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
)

// deleteCmd represents the add command
var deleteCmd = &cobra.Command{
	Use:   "delete [service]",
	Short: "delete the given service from crypt",
	Long: `Delete a service from the crypt getStore(). Note that the
deleted service cannot be recovered.`,
	Args:    serviceIsValid,
	RunE:    delete,
	Aliases: []string{"del", "remove"},
}

func delete(cmd *cobra.Command, args []string) error {
	asker := asker.DefaultAsker()
	ok, err := asker.AskConfirm("Are you sure you want delete this service?")
	if err != nil {
		return err
	}
	if ok {
		service := args[0]
		st, err := getStore()
		if err != nil {
			return err
		}
		st.RemoveCredential(service)
		color.Blue("successfully deleted: %s", service)
		saveStore()
	}
	return nil
}
