package cmds

import (
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/creds"
	"github.com/sugatpoudel/crypt/internal/utils"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [service]",
	Short: "Add a service to crypt",
	Long: `Add a service along with any associated information
to the crypt getStore().

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

	email, err := asker.Ask("Email")
	utils.FatalIf(err)

	user, err := asker.Ask("Username")
	utils.FatalIf(err)

	pwd, err := asker.AskSecret("Password", true)
	utils.FatalIf(err)

	desc, err := asker.Ask("Description")
	utils.FatalIf(err)

	cred := creds.Credential{
		Service:     service,
		Email:       email,
		Username:    user,
		Password:    pwd,
		Description: desc,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	cred.PrintCredential()

	ok, err := asker.AskConfirm("Does this look right?")
	if ok {
		getStore().Crypt.SetCredential(cred)
		color.Green("\nAdded service '%s'", service)
	}
	color.Green("Saving crypt")

	err = getStore().Save()
	utils.FatalIf(err)
}
