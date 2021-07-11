package cobracli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/crypt"

	"github.com/spf13/cobra"
)

var (
	fields = []string{"email", "username", "password", "description", "tags", "exit"}
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [service]",
	Short: "edit fields for the given service",
	Long: `edit fields for the given service.
Similar flow to the add command however, blank
values are interpreted as a no-op.
`,
	Args: combineArgs(backupCrypt, parseService),
	RunE: edit,
}

func edit(cmd *cobra.Command, args []string) error {
	svc, err := getService()
	if err != nil {
		return err
	}

	asker := asker.DefaultAsker()
	updated := &crypt.Credential{Service: svc.Service}
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
			val, err := asker.Ask("tag")
			if err != nil {
				return err
			}
			tags := strings.Split(val, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
			updated.Tags = append(updated.Tags, tags...)
		case 5:
			exit = true
		}

		if exit {
			break
		}
	}

	crypt.PrintCredential(updated)
	fmt.Println()

	msg := color.YellowString("do these updated values make sense?")
	ok, err := asker.AskConfirm(msg)
	if err != nil {
		return err
	}
	if ok {
		svc.Merge(updated)
		color.Green("updated service '%s'", svc.Service)
		saveStore()
	}
	return nil
}
