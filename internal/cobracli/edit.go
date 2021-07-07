package cobracli

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/creds"
	"github.com/tagus/crypt/internal/finder"
	"github.com/tagus/crypt/internal/utils"

	"github.com/spf13/cobra"
)

var (
	fields = []string{"email", "username", "password", "description", "exit"}
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [service]",
	Short: "edit fields for the given service",
	Long: `edit fields for the given service.
Similar flow to the add command however, blank
values are interpreted as a no-op.
`,
	Args: serviceIsValid,
	RunE: edit,
}

func edit(cmd *cobra.Command, args []string) error {
	service := args[0]
	asker := asker.DefaultAsker()

	st, err := getStore()
	if err != nil {
		return err
	}
	fd := finder.New(st.Crypt)
	oldCred := fd.Find(service)

	var email, user, pwd, desc string
	for {
		n, err := asker.AskSelect("what would you like to edit?", fields)
		if err != nil {
			return err
		}

		exit := false
		switch n {
		case 0:
			email, err = asker.Ask("email")
			if err != nil {
				return err
			}
		case 1:
			user, err = asker.Ask("username")
			if err != nil {
				return err
			}
		case 2:
			pwd, err = asker.AskSecret("pwd", true)
			if err != nil {
				return err
			}
		case 3:
			desc, err = asker.Ask("description")
			if err != nil {
				return err
			}
		case 4:
			exit = true
		}

		if exit {
			break
		}
	}

	cred := &creds.Credential{
		Service:     oldCred.Service,
		Email:       utils.FallbackStr(email, oldCred.Email),
		Username:    utils.FallbackStr(user, oldCred.Username),
		Password:    utils.FallbackStr(pwd, oldCred.Password),
		Description: utils.FallbackStr(desc, oldCred.Description),
		CreatedAt:   oldCred.CreatedAt,
		UpdatedAt:   time.Now().Unix(),
	}

	creds.PrintCredential(cred)
	fmt.Println()

	msg := color.YellowString("does this look right?")
	ok, err := asker.AskConfirm(msg)
	utils.FatalIf(err)
	if ok {
		_, err := st.SetCredential(cred)
		if err != nil {
			return err
		}
		color.Green("updated service '%s'", service)
		saveStore()
	}
	return nil
}
