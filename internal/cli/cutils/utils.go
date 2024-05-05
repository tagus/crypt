package cutils

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/repos"
	"strings"
)

// ServiceIsNew ensures that the given service does not already exist
func ServiceIsNew(cmd *cobra.Command, args []string) error {
	svc, err := ParseService(cmd, args)
	if err != nil {
		return err
	}

	env, err := environment.Load(cmd)
	if err != nil {
		return err
	}

	repo := env.Repo()
	creds, err := repo.QueryCredentials(cmd.Context(), repos.QueryCredentialsFilter{Service: svc})
	if err != nil {
		return err
	}

	if len(creds) > 0 {
		return fmt.Errorf("service %s already exists", svc)
	}
	return nil
}

// ParseService parses the args to a single service string
func ParseService(cmd *cobra.Command, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("no service provided")
	}
	combined := strings.Join(args, " ")
	combined = strings.TrimSpace(combined)
	return combined, nil
}

func CombineArgs(fns ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range fns {
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
