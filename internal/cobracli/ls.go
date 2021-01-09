// Copyright © 2018 Sugat Poudel <taguspoudel@gmail.com>

package cobracli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/utils"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "list stored services",
	Long:    `lists the name of all stored service credentials.`,
	RunE:    ls,
	Aliases: []string{"list"},
}

func ls(cmd *cobra.Command, args []string) error {
	st, err := getStore()
	if err != nil {
		return err
	}

	creds := st.Credentials

	data := make([][]string, len(creds))
	counter := 0
	for _, v := range creds {
		createdAt := v.GetCreatedAt().Format("01/02/2006")
		data[counter] = []string{strconv.Itoa(counter), v.Service, createdAt}
		counter++
	}

	utils.PrintTable(data, utils.TableOpts{
		Headers: []string{"index", "name", "created at"},
	})
	fmt.Printf("%d credential(s).\n", len(creds))

	return nil
}
