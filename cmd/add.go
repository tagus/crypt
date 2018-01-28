package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/backend"
	input "github.com/tcnksm/go-input"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [service_name]",
	Short: "Add credentials for a service",
	Long: `Add the complete credentials for a service.

The first argument is the name of the service`,
	Args: cobra.ExactArgs(1),
	Run:  addService,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addService(cmd *cobra.Command, args []string) {
	service := args[0]
	ui := input.DefaultUI()

	desc, err := ui.Ask("Describe this service.", &input.Options{
		Required: true,
		Loop:     true,
	})
	backend.PrintErrorAndExit(err)

	email, err := ui.Ask("Email?", &input.Options{
		Required: true,
		Loop:     true,
	})
	backend.PrintErrorAndExit(err)

	pwd, err := ui.Ask("Password?", &input.Options{
		Required: true,
		Loop:     true,
	})
	backend.PrintErrorAndExit(err)

	cred := backend.Credential{
		Service:     service,
		Email:       email,
		Password:    pwd,
		Description: desc,
	}

	CryptFile.AddCredential(cred)
	err = backend.ReEncode(CryptFile, CryptPath, KeyString)
	backend.PrintErrorAndExit(err)

	color.Green("%s was successfully added", service)
}
