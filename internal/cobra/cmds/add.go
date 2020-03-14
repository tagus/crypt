package cmds

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/creds"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [service]",
	Short: "Add a sevice to crypt",
	Long: `Add a service along with any associated information
to the crypt store.

Expects a single argument, however multi word services
can be espaced using quotes.`,
	Args:    serviceIsNew,
	Example: "add 'Amazon Web Services'",
	Run:     add,
	Aliases: []string{"new"},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func add(cmd *cobra.Command, args []string) {
	service := args[0]
	asker := asker.DefaultAsker()

	email, err := asker.Ask("Email: ", nil)
	printAndExit(err)

	user, err := asker.Ask("Username: ", nil)
	printAndExit(err)

	pwd, err := asker.AskSecret("Password: ", true, nil)
	printAndExit(err)

	desc, err := asker.Ask("Description: ", nil)
	printAndExit(err)

	cred := creds.Credential{
		Service:     service,
		Email:       email,
		Username:    user,
		Password:    pwd,
		Description: desc,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	fmt.Println()
	cred.PrintCredential()

	msg := color.YellowString("\nDoes this look right? [y/n]")
	confirm, err := asker.Ask(msg, nil)
	if confirm == "y" {
		Store.Crypt.SetCredential(cred)
		color.Green("Added service '%s'", service)
	}
	color.Green("\nSaving crypt")

	err = Store.Save()
	printAndExit(err)
}
