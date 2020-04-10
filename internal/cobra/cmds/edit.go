package cmds

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sugatpoudel/crypt/internal/asker"
	"github.com/sugatpoudel/crypt/internal/creds"
	"github.com/sugatpoudel/crypt/internal/utils"

	"github.com/spf13/cobra"
)

var (
	fields = []string{"email", "username", "password", "description", "exit"}
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [service]",
	Short: "Edit fields for the given service",
	Long: `Edit fields for the given service.
Similar flow to the add command however, blank
values are interpreted as a no-op.
`,
	Args: serviceIsValid,
	Run:  edit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func edit(cmd *cobra.Command, args []string) {
	service := args[0]
	asker := asker.DefaultAsker()

	oldCred := getStore().Crypt.FindCredential(service)

	var email, user, pwd, desc string
	for {
		n, err := asker.AskSelect("What would you like to edit?", fields)
		utils.FatalIf(err)

		exit := false
		switch n {
		case 0:
			email, err = asker.Ask("Email")
			utils.FatalIf(err)
		case 1:
			user, err = asker.Ask("Username")
			utils.FatalIf(err)
		case 2:
			pwd, err = asker.AskSecret("Password", true)
			utils.FatalIf(err)
		case 3:
			desc, err = asker.Ask("Description")
			utils.FatalIf(err)
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

	cred.PrintCredential()
	fmt.Println()

	msg := color.YellowString("Does this look right?")
	ok, err := asker.AskConfirm(msg)
	utils.FatalIf(err)
	if ok {
		getStore().Crypt.SetCredential(cred)
		color.Green("Updated service '%s'", service)
		saveStore()
	}
}

// Returns the old string if the new string is empty
func noop(old, new string) string {
	if strings.TrimSpace(new) == "" {
		return old
	}
	return new
}
