// Copyright © 2018 Sugat Poudel <taguspoudel@gmail.com>

package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/utils"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List stored services",
	Long:    `Lists the name of all stored service credentials.`,
	Run:     ls,
	Aliases: []string{"list"},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func ls(cmd *cobra.Command, args []string) {
	creds := Store.Crypt.Credentials

	data := make([][]string, len(creds))
	counter := 0
	for _, v := range creds {
		createdAt := v.GetCreatedAt().Format("Jan _2 2006")
		data = append(data, []string{strconv.Itoa(counter), v.Service, createdAt})
		counter++
	}

	caption := fmt.Sprintf("%d credential(s).", len(creds))
	utils.PrintTable(data, &caption)
	fmt.Println()
}
