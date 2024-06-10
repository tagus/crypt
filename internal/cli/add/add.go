package add

import (
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/cli/cutils"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/components"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/crypt/internal/utils"
	"github.com/tagus/mango"
)

var (
	useUI bool
)

var Command = &cobra.Command{
	Use:     "add [service]",
	Short:   "add a service to crypt",
	Long:    `add a service along with any associated information to the crypt`,
	Args:    cutils.ServiceMaybeNew,
	Example: "add Amazon Web Services",
	RunE:    add,
}

func init() {
	Command.Flags().BoolVar(&useUI, "ui", false, "whether to use the ui form view instead")
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
	if useUI {
		form := components.NewForm(components.FormOpts{
			Title: "add a new service credential",
			OnSave: func(cr *repos.Credential, fn components.ShowModalFn) (*repos.Credential, error) {
				return env.Repo().InsertCredential(cmd.Context(), cr)
			},
		})
		cred, err = form.Show()
		if err != nil {
			return err
		}
	} else {
		cred, err = getCredentialDetails(svc)
		if err != nil {
			return err
		}
		repo := env.Repo()
		cred, err = repo.InsertCredential(cmd.Context(), cred)
		if err != nil {
			return err
		}
	}

	if cred == nil {
		mango.Warning("no credential was added")
		return nil
	}
	mango.Debug(cred)
	mango.Debug("added credential for", cred.Service)
	return nil
}

/******************************************************************************/

func getCredentialDetails(service string) (*repos.Credential, error) {
	var (
		cred *repos.Credential
		err  error
	)
	for {
		cred, err = buildCredential(service)
		if err != nil {
			return nil, err
		}
		if cred != nil {
			break
		}
	}
	return cred, nil
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
