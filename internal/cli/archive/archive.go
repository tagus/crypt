package archive

import (
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/cli/cutils"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/mango"
)

var Command = &cobra.Command{
	Use:     "archive [service]",
	Short:   "deletes the given service from the crypt db",
	RunE:    archive,
	Aliases: []string{"delete", "rm"},
}

func archive(cmd *cobra.Command, args []string) error {
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
	if err := repo.ArchiveCredential(cmd.Context(), credID); err != nil {
		return err
	}
	mango.Debug("archived credential:", credID)

	return nil
}
