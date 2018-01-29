package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [service]",
	Short: "Add a sevice to crypt.",
	Long: `Add a service along with any associated information
to the crypt store.

Expects a single argument, however multi word services
can be espaced using quotes.`,
	Args:    cobra.ExactArgs(1),
	Example: "add 'Amazon Web Services'",
	Run:     add,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func add(cmd *cobra.Command, args []string) {

}
