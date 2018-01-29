package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
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
	var email, user, desc string

	fmt.Print("Email: ")
	fmt.Scanln(&email)

	fmt.Print("Username: ")
	fmt.Scanln(&user)

	fmt.Print("Password: ")
	pwdB, err := gopass.GetPasswd()
	printAndExit(err)
	pwd := string(pwdB)

	fmt.Print("Decription: ")
	fmt.Scanln(&desc)

	fmt.Printf("\nEmail: %s\n", email)
	fmt.Printf("Username: %s\n", user)
	fmt.Printf("Password: [redacted]\n")
	fmt.Printf("Description: %s\n", desc)

	msg := color.YellowString("Does this look right? [y/n]")
	fmt.Printf("%s ", msg)
	var confirm string
	fmt.Scanln(&confirm)

	if confirm == "y" {
		cred := creds.Credential{
			Service:     service,
			Email:       email,
			Username:    user,
			Password:    pwd,
			Description: desc,
		}
		Store.Crypt.SetCredential(cred)
		color.Green("Added service '%s'", service)
	}

	color.Green("\nSaving crypt")
	err = Store.Save()
	printAndExit(err)
}
