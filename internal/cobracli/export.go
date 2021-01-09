package cobracli

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
	Short: "export the cryptfile to plain json",
	Long: `export will decrypt the current cryptfile and export
all credentials as plain text. This is purely meant as a convenience
function and should be used sparingly.`,
	RunE: export,
	Args: cobra.ExactArgs(1),
}

func export(cmd *cobra.Command, args []string) error {
	asker := asker.DefaultAsker()
	ok, err := asker.AskConfirm("are you sure you want to export the cryptfile?")
	utils.FatalIf(err)

	if ok {
		st, err := getStore()
		if err != nil {
			return err
		}

		fmt.Println("exporting cryptfile")
		data, err := st.GetJSON()
		if err != nil {
			return err
		}

		output := args[0]
		err = ioutil.WriteFile(output, data, 0644)
		if err != nil {
			return err
		}

		fmt.Println("cryptfile exported to: ", output)
	}
	return err
}
