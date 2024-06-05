package edit

import (
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/cli/cutils"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/utils"
	"github.com/tagus/mango"
)

var Command = &cobra.Command{
	Use:   "edit [service]",
	Short: "edit fields for the given service",
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	env, err := environment.Load(cmd)
	if err != nil {
		return err
	}
	defer env.Close()

	credID, err := cutils.GetCredentialID(env, cmd, args)
	if err != nil {
		return err
	}
	mango.Debug("selected credential id:", credID)

	repo := env.Repo()
	cred, err := repo.AccessCredential(cmd.Context(), credID)
	if err != nil {
		return err
	}

	mango.Debug("editing credential:", cred.Service)

	form := &Form{cr: env.Repo()}
	cred, err = form.Show(cmd.Context(), cred)
	if err != nil {
		return err
	}

	if cred == nil {
		mango.Warning("no credential was updated")
		return nil
	}
	mango.Debug("updated credential:", cred.Service)
	utils.PrintCredential(cred)

	return nil
}
