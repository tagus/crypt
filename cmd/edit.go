package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/asker"
	"github.com/sugatpoudel/crypt/creds"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [service]",
	Short: "Edit fields for the given service",
	Long: `Edit fields for the given service.
Similar flow to the add command however, blank
values are interpreted as a no-op.`,
	Args: serviceIsValid,
	Run:  edit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

// Returns the old string if the new string is empty
func noop(old, new string) string {
	if strings.TrimSpace(new) == "" {
		return old
	}
	return new
}

func edit(cmd *cobra.Command, args []string) {
	service := args[0]
	asker := asker.DefaultAsker()

	oldCred := Store.Crypt.FindCredential(service)

	email, err := asker.Ask("Email: ", nil)
	printAndExit(err)

	user, err := asker.Ask("Username: ", nil)
	printAndExit(err)

	pwd, err := asker.AskSecret("Password: ", true, nil)
	printAndExit(err)

	desc, err := asker.Ask("Description: ", nil)
	printAndExit(err)

	cred := creds.Credential{
		Service:     oldCred.Service,
		Email:       noop(oldCred.Email, email),
		Username:    noop(oldCred.Username, user),
		Password:    noop(oldCred.Password, pwd),
		Description: noop(oldCred.Description, desc),
	}

	fmt.Println()
	cred.PrintCredential()

	msg := color.YellowString("\nDoes this look right? [y/n]")
	confirm, err := asker.Ask(msg, nil)
	if confirm == "y" {
		Store.Crypt.SetCredential(cred)
		color.Green("Updated service '%s'", service)
	}
	color.Green("\nSaving crypt")

	err = Store.Save()
	printAndExit(err)
}
