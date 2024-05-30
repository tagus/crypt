package show

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/cli/cutils"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/crypt/internal/utils"
	"github.com/tagus/mango"
)

var Command = &cobra.Command{
	Use:   "show",
	Short: "show information about a service",
	Long:  `prompts user for a service and shows the results in a selectable list`,
	RunE:  show,
}

func show(cmd *cobra.Command, args []string) error {
	env, err := environment.Load(cmd)
	if err != nil {
		return err
	}
	defer env.Close()

	svc, err := cutils.ParseService(cmd, args)
	if err != nil {
		return err
	}
	mango.Debug("show query:", svc)

	var (
		crypt = env.Crypt()
		repo  = env.Repo()
	)
	creds, err := repo.QueryCredentials(cmd.Context(), repos.QueryCredentialsFilter{Service: svc, CryptID: crypt.ID})
	if err != nil {
		return err
	}

	if len(creds) == 0 {
		return fmt.Errorf("no credential found matching query: %s", svc)
	}

	cred := creds[0]
	if len(creds) > 1 {
		cred, err = selectCredential(creds)
		if err != nil {
			return err
		}
	}

	cred, err = repo.AccessCredential(cmd.Context(), crypt.ID, cred.ID)
	if err != nil {
		return err
	}

	utils.PrintCredential(cred)

	return nil
}

func selectCredential(creds []*repos.Credential) (*repos.Credential, error) {
	ak := asker.DefaultAsker()
	i, err := ak.AskSelect("select service", creds)
	if err != nil {
		return nil, err
	}
	return creds[i], nil
}
