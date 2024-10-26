package edit

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/cli/cutils"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/utils"
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
	slog.Debug("selected credential", "id", credID)

	repo := env.Repo()
	cred, err := repo.AccessCredential(cmd.Context(), credID)
	if err != nil {
		return err
	}

	slog.Debug("editing credential", "service", cred.Service)

	form := &Form{cr: env.Repo()}
	cred, err = form.Show(cmd.Context(), cred)
	if err != nil {
		return err
	}

	if cred == nil {
		slog.Warn("no credential was updated")
		return nil
	}
	slog.Debug("updated credential", "service", cred.Service)
	utils.PrintCredential(cred)

	return nil
}
