package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/utils"
)

// infoCmd represents the export command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about your crypt file",
	Long:  `Displays meta information about your crypt file.`,
	Run:   info,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) {
	crypt := Store.Crypt
	title := ` _______  ______    __   __  _______  _______ 
|       ||    _ |  |  | |  ||       ||       |
|       ||   | ||  |  |_|  ||    _  ||_     _|
|       ||   |_||_ |       ||   |_| |  |   |  
|      _||    __  ||_     _||    ___|  |   |  
|     |_ |   |  | |  |   |  |   |      |   |  
|_______||___|  |_|  |___|  |___|      |___|  
	`

	data := [][]string{
		[]string{"credentials", strconv.Itoa(len(crypt.Credentials))},
		[]string{"created at", crypt.GetCreatedAt().Format("Jan _2 2006")},
		[]string{"created at", crypt.GetUpdatedAt().Format("Jan _2 2006")},
	}

	fmt.Println(title)
	utils.PrintTable(data, nil)
}
