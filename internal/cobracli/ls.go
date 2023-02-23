package cobracli

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/crypt"
	"github.com/tagus/crypt/internal/utils"
)

var lsLimit int

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "list stored services",
	Long:    `lists the name of all stored service credentials.`,
	RunE:    ls,
	Aliases: []string{"list"},
}

func init() {
	lsCmd.Flags().IntVarP(&lsLimit, "limit", "l", 0, "limit the number of services to list")
}

func ls(cmd *cobra.Command, args []string) error {
	st, err := getStore()
	if err != nil {
		return err
	}

	creds := make([]*crypt.Credential, 0, len(st.Credentials))
	for _, v := range st.Credentials {
		creds = append(creds, v)
	}
	sort.Slice(creds, func(i, j int) bool {
		if creds[i].AccessedAt == nil || creds[j].AccessedAt == nil {
			return creds[j].AccessedAt == nil
		}
		return *creds[i].AccessedAt > *creds[j].AccessedAt
	})

	if lsLimit > 0 {
		creds = creds[:utils.Min(lsLimit, len(creds))]
	}

	data := make([][]string, len(creds))
	counter := 0
	for _, v := range creds {
		data[counter] = []string{
			strconv.Itoa(counter),
			v.Service,
			utils.FormatTimeSince(v.GetAccessedAt()),
		}
		counter++
	}

	utils.PrintTable(data, utils.TableOpts{
		Headers: []string{"index", "name", "last accessed at"},
	})
	fmt.Printf("%d total credential(s).\n", len(st.Credentials))

	return nil
}
