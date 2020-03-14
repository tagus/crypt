package cmds

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/creds"
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

	isNumber := func(str string) error {
		_, err := strconv.Atoi(strings.Trim(str, " \n"))
		if err != nil {
			return err
		}
		return nil
	}

	var email, user, pwd, desc string
	for {
		ans, err := asker.Ask("What would you like to edit? ", isNumber)
		printAndExit(err)
		exit := false

		switch i, _ := strconv.Atoi(strings.Trim(ans, " ")); i {
		case 1:
			email, err = asker.Ask("Email: ")
			printAndExit(err)
		case 2:
			user, err = asker.Ask("Username: ")
			printAndExit(err)
		case 3:
			pwd, err = asker.AskSecret("Password: ", true)
			printAndExit(err)
		case 4:
			desc, err = asker.Ask("Description: ")
			printAndExit(err)
		default:
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
