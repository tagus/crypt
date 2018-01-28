package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all stored services",
	Long: `Lists the names of all stored services
as stored in the crypt file.

Also prints meta information regarding sotred services.`,
	Args: cobra.ExactArgs(0),
	Run:  listServices,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listServices(cmd *cobra.Command, args []string) {
	credentials := CryptFile.Credentials
	fmt.Printf("%s total services\n", color.YellowString("%d", len(credentials)))
	for _, c := range credentials {
		fmt.Println(c.Service)
	}
}
