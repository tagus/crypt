// Copyright © 2018 Sugat Poudel <taguspoudel@gmail.com>

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List stored services",
	Long:  `Lists the name of all stored service credentials.`,
	Run:   ls,
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func ls(cmd *cobra.Command, args []string) {
	creds := Store.Crypt.Credentials
	fmt.Printf("You have %d stored credential(s).\n", len(creds))
	for _, v := range creds {
		fmt.Println(v.Service)
	}
}
