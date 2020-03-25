package cmds

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/creds"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [service]",
	Short: "Edit fields for the given service",
	Long: `Edit fields for the given service.
Similar flow to the add command however, blank
values are interpreted as a no-op.

The following are the numerical values for each field

	1. email
	2. username
	3. password
	4. description

`,
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

	var email, user, pwd, desc string
	for {
		prompt := promptui.Select{
			Label: "What would you like to edit?",
			Items: []string{"email", "username", "password", "description", "exit"},
		}
		n, _, err := prompt.Run()
		printAndExit(err)

		exit := false

		switch n {
		case 0:
			email, err = asker.Ask("Email: ")
			printAndExit(err)
		case 1:
			user, err = asker.Ask("Username: ")
			printAndExit(err)
		case 2:
			pwd, err = asker.AskSecret("Password: ", true)
			printAndExit(err)
		case 3:
			desc, err = asker.Ask("Description: ")
			printAndExit(err)
		case 4:
			exit = true
		}

		if exit {
			break
		}
	}

	cred := creds.Credential{
		Service:     oldCred.Service,
		Email:       noop(oldCred.Email, email),
		Username:    noop(oldCred.Username, user),
		Password:    noop(oldCred.Password, pwd),
		Description: noop(oldCred.Description, desc),
		CreatedAt:   oldCred.CreatedAt,
		UpdatedAt:   time.Now().Unix(),
	}

	fmt.Println()
	cred.PrintCredential()

	msg := color.YellowString("\nDoes this look right? [y/n]")
	confirm, err := asker.Ask(msg)
	if confirm == "y" {
		Store.Crypt.SetCredential(cred)
		color.Green("Updated service '%s'", service)
	}
	color.Green("\nSaving crypt")

	err = Store.Save()
	printAndExit(err)
}
