package cmds

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/utils"
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
	asker := asker.DefaultAsker()
	ok, err := asker.AskConfirm("Are you sure you want to export the cryptfile?")
	utils.FatalIf(err)

	if ok {
		fmt.Println("exporting cryptfile")
		data, err := getStore().Crypt.GetJSON()
		utils.FatalIf(err)

		output := args[0]
		err = ioutil.WriteFile(output, data, 0644)
		utils.FatalIf(err)

		fmt.Println("cryptfile exported to: ", output)
	}
}
