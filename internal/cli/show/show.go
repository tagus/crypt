package show

import (
	"errors"
	"log/slog"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/cli/cutils"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/utils"
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

	credID, err := cutils.GetCredentialID(env, cmd, args)
	if err != nil {
		return err
	}
	slog.Debug("selected credential id", "cred_id", credID)

	repo := env.Repo()
	cred, err := repo.AccessCredential(cmd.Context(), credID)
	if err != nil {
		return err
	}

	utils.PrintCredential(cred)

	if err := clipboard.WriteAll(cred.Password); err != nil {
		return errors.New("failed to copy pwd")
	}
	slog.Debug("copied password to clipboard")

	return nil
}
