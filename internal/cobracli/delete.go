package cobracli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
)

// deleteCmd represents the add command
var deleteCmd = &cobra.Command{
	Use:   "delete [service]",
	Short: "delete the given service from crypt",
	Long: `Delete a service from the crypt getStore(). Note that the
deleted service cannot be recovered.`,
	Args:    combineArgs(backupCrypt, parseService),
	RunE:    delete,
	Aliases: []string{"del", "rm"},
}

func delete(cmd *cobra.Command, args []string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	asker := asker.DefaultAsker()
	ok, err := asker.AskConfirm(fmt.Sprintf("are you sure you want delete %s?", svc.Service))
	if err != nil {
		return err
	}
	if ok {

		st, err := getStore()
		if err != nil {
			return err
		}
		st.RemoveCredential(svc.Id)
		color.Blue("successfully deleted: %s", svc.Service)
		saveStore()
	}
	return nil
}
