package cmds

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/utils"
)

// deleteCmd represents the add command
var deleteCmd = &cobra.Command{
	Use:   "delete [service]",
	Short: "delete the given service from crypt",
	Long: `Delete a service from the crypt getStore(). Note that the
deleted service cannot be recovered.`,
	Args:    serviceIsValid,
	Run:     delete,
	Aliases: []string{"del", "remove"},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func delete(cmd *cobra.Command, args []string) {
	asker := asker.DefaultAsker()
	ok, err := asker.AskConfirm("Are you sure you want delete this service?")
	utils.FatalIf(err)
	if ok {
		service := args[0]
		getStore().Crypt.RemoveCredential(service)
		color.Blue("successfully deleted: %s", service)
		saveStore()
	}
}
