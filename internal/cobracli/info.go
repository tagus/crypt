package cobracli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sugatpoudel/crypt/internal/utils"
)

// infoCmd represents the export command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about your crypt file",
	Long:  `Displays meta information about your crypt file.`,
	RunE:  info,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func info(cmd *cobra.Command, args []string) error {
	st, err := getStore()
	if err != nil {
		return err
	}
	title := ` _______  ______    __   __  _______  _______
|       ||    _ |  |  | |  ||       ||       |
|       ||   | ||  |  |_|  ||    _  ||_     _|
|       ||   |_||_ |       ||   |_| |  |   |
|      _||    __  ||_     _||    ___|  |   |
|     |_ |   |  | |  |   |  |   |      |   |
|_______||___|  |_|  |___|  |___|      |___|
	`

	data := [][]string{
		{"credentials", strconv.Itoa(st.Len())},
		{"created at", st.GetCreatedAt().Format("Jan _2 2006")},
		{"updated at", st.GetUpdatedAt().Format("Jan _2 2006")},
	}

	fmt.Println(title)
	utils.PrintTable(data, utils.TableOpts{})
	return nil
}
