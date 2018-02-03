package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/asker"
	"github.com/sugatpoudel/crypt/creds"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [service]",
	Short: "Add a sevice to crypt.",
	Long: `Add a service along with any associated information
to the crypt store.

Expects a single argument, however multi word services
can be espaced using quotes.`,
	Args:    cobra.ExactArgs(1),
	Example: "add 'Amazon Web Services'",
	Run:     add,
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

	fmt.Print("Password: ")
	pwdB, err := gopass.GetPasswd()
	printAndExit(err)
	pwd := string(pwdB)

	fmt.Print("Confirm Password: ")
	pwdB2, err := gopass.GetPasswd()
	printAndExit(err)
	pwd2 := string(pwdB2)

	if pwd == "" || pwd != pwd2 {
		color.Red("Passwords did not match. Try again.")
		return
	}

	desc, err := asker.Ask("Description: ", nil)
	printAndExit(err)

	cred := creds.Credential{
		Service:     service,
		Email:       email,
		Username:    user,
		Password:    pwd,
		Description: desc,
	}

	fmt.Println()
	cred.PrintCredential()
	msg := color.YellowString("\nDoes this look right? [y/n]")
	fmt.Printf("%s ", msg)
	var confirm string
	fmt.Scanln(&confirm)

	if confirm == "y" {
		Store.Crypt.SetCredential(cred)
		color.Green("Added service '%s'", service)
	}
	color.Green("\nSaving crypt")

	err = Store.Save()
	printAndExit(err)
}
