package list

import (
	"fmt"
	"github.com/tagus/mango"

	"github.com/tagus/crypt/internal/repos"

	"github.com/tagus/crypt/internal/cli/environment"

	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/utils"
)

var (
	limit     int
	tagFilter string
)

var Command = &cobra.Command{
	Use:     "list",
	Short:   "list stored credentials",
	Long:    `lists the name of all stored service credentials.`,
	RunE:    list,
	Args:    cobra.NoArgs,
	Aliases: []string{"ls"},
}

func init() {
	Command.Flags().IntVarP(&limit, "limit", "l", 0, "limit the number of services to list")
	Command.Flags().StringVarP(&tagFilter, "tag", "t", "", "filter services by tag")
}

func list(cmd *cobra.Command, args []string) error {
	env, err := environment.Load(cmd)
	if err != nil {
		return err
	}
	defer env.Close()

	repo := env.Repo()
	creds, err := repo.QueryCredentials(cmd.Context(), repos.QueryCredentialsFilter{
		Limit: limit,
		Tag:   tagFilter,
	})

	if len(creds) == 0 {
		fmt.Println("no credentials found")
		return nil
	}

	data := make([][]string, len(creds))
	counter := 0
	for _, cred := range creds {
		data[counter] = []string{
			cred.ID,
			cred.Service,
			mango.FormatTimeSince(cred.AccessedAt),
		}
		counter++
	}

	fmt.Printf("%d total credential(s).\n", len(creds))
	utils.PrintTable(data, utils.TableOpts{
		Headers: []string{"id", "name", "last accessed at"},
	})

	return nil
}
