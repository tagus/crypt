package cmds

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
)

// deleteCmd represents the add command
var deleteCmd = &cobra.Command{
	Use:   "delete [service]",
	Short: "delete the given service from crypt",
	Long: `Delete a service from the crypt store. Note that the
deleted service cannot be recovered.`,
	Args:    serviceIsValid,
	Run:     delete,
	Aliases: []string{"new"},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func delete(cmd *cobra.Command, args []string) {
	ak := asker.DefaultAsker()
	confirm, err := ak.Ask("Are you sure you want delete this service?")
	printAndExit(err)

	if confirm == "yes" {
		service := args[0]
		Store.Crypt.RemoveCredential(service)

		color.Blue("successfully deleted: %s", service)
		color.Green("\nSaving crypt")

		err := Store.Save()
		printAndExit(err)
	}
}
