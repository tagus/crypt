package add

import (
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/cli/cutils"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/crypt/internal/utils"
	"github.com/tagus/mango"
)

var Command = &cobra.Command{
	Use:     "add [service]",
	Short:   "add a service to crypt",
	Long:    `add a service along with any associated information to the crypt`,
	Args:    cutils.ServiceIsNew,
	Example: "add Amazon Web Services",
	RunE:    add,
}

func add(cmd *cobra.Command, args []string) error {
	env, err := environment.Load(cmd)
	if err != nil {
		return err
	}
	defer env.Close()

	svc, err := cutils.ParseService(cmd, args)
	if err != nil {
		return err
	}

	var cred *repos.Credential
	for {
		cred, err = buildCredential(svc)
		if err != nil {
			return err
		}
		if cred != nil {
			break
		}
	}

	crypt := env.Crypt()
	repo := env.Repo()
	cred, err = repo.InsertCredential(cmd.Context(), crypt.ID, cred)
	if err != nil {
		return err
	}
	mango.Info(cred)
	mango.Info("added credential for", cred.Service)
	return nil
}

func buildCredential(service string) (*repos.Credential, error) {
	ak := asker.DefaultAsker()

	email, err := ak.Ask("email")
	if err != nil {
		return nil, err
	}

	user, err := ak.Ask("username")
	if err != nil {
		return nil, err
	}

	pwd, err := ak.AskSecret("pwd", true)
	if err != nil {
		return nil, err
	}

	desc, err := ak.Ask("description")
	if err != nil {
		return nil, err
	}

	// TODO: ask for other metadata

	cred := &repos.Credential{
		ID:          mango.ShortID(),
		Service:     service,
		Email:       email,
		Username:    user,
		Password:    pwd,
		Description: desc,
		Details:     &repos.Details{},
		Tags:        []string{},
		Domains:     []string{},
	}
	utils.PrintCredential(cred)

	ok, err := ak.AskConfirm("does this look right?")
	if !ok {
		return nil, nil
	}

	return cred, nil
}
