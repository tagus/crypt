package cutils

import (
	"errors"

	"github.com/tagus/crypt/internal/asker"

	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/repos"
)

var (
	ErrNoCredentialFound = errors.New("no credentials found")
)

func GetCredentialID(
	env *environment.Environment,
	cmd *cobra.Command,
	args []string,
) (string, error) {
	svc, err := ParseService(cmd, args)
	if err != nil {
		return "", err
	}

	repo := env.Repo()
	creds, err := repo.QueryCredentials(cmd.Context(), repos.QueryCredentialsFilter{Service: svc})
	if err != nil {
		return "", err
	}

	if len(creds) == 0 {
		return "", ErrNoCredentialFound
	}

	if len(creds) > 1 {
		cred, err := selectCredential(creds)
		if err != nil {
			return "", err
		}
		return cred.ID, nil
	}
	return creds[0].ID, nil
}

func selectCredential(creds []*repos.Credential) (*repos.Credential, error) {
	ak := asker.DefaultAsker()
	i, err := ak.AskSelect("select service", creds)
	if err != nil {
		return nil, err
	}
	return creds[i], nil
}
