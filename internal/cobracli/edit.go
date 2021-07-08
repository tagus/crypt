package cobracli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/creds"
	"github.com/tagus/crypt/internal/utils"

	"github.com/spf13/cobra"
)

var (
	fields = []string{"email", "username", "password", "description", "exit"}
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [service]",
	Short: "edit fields for the given service",
	Long: `edit fields for the given service.
Similar flow to the add command however, blank
values are interpreted as a no-op.
`,
	Args: parseService,
	RunE: edit,
}

func edit(cmd *cobra.Command, args []string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	st, err := getStore()
	if err != nil {
		return err
	}

	asker := asker.DefaultAsker()
	updated := &creds.Credential{}
	for {
		n, err := asker.AskSelect("what would you like to edit?", fields)
		if err != nil {
			return err
		}

		exit := false
		switch n {
		case 0:
			updated.Email, err = asker.Ask("email")
			if err != nil {
				return err
			}
		case 1:
			updated.Username, err = asker.Ask("username")
			if err != nil {
				return err
			}
		case 2:
			updated.Password, err = asker.AskSecret("pwd", true)
			if err != nil {
				return err
			}
		case 3:
			updated.Description, err = asker.Ask("description")
			if err != nil {
				return err
			}
		case 4:
			exit = true
		}

		if exit {
			break
		}
	}

	creds.PrintCredential(updated)
	fmt.Println()

	msg := color.YellowString("do these updated values make sense?")
	ok, err := asker.AskConfirm(msg)
	utils.FatalIf(err)
	if ok {
		_, err := st.SetCredential(svc.Merge(updated))
		if err != nil {
			return err
		}
		color.Green("updated service '%s'", updated.Service)
		saveStore()
	}
	return nil
}
