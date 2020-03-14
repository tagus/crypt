package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/sugatpoudel/crypt/asker"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export [output]",
	Short: "Export the cryptfile to plain json",
	Long: `Export will decrypt the current cryptfile and export
all credentials as plain text. This is purely meant as a convenience
function and should be used sparingly.`,
	Run:  export,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func export(cmd *cobra.Command, args []string) {
	ak := asker.DefaultAsker()
	confirm, err := ak.Ask("Are you sure you want to export the cryptfile?", nil)
	printAndExit(err)

	if confirm == "yes" {
		fmt.Println("exporting cryptfile")
		data, err := Store.Crypt.GetJSON()
		printAndExit(err)

		output := args[0]
		err = ioutil.WriteFile(output, data, 0644)
		printAndExit(err)

		fmt.Println("cryptfile exported to: ", output)
	}
}
