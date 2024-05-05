package info

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/utils"
	"strconv"
)

var title = ` _______  ______    __   __  _______  _______
|       ||    _ |  |  | |  ||       ||       |
|       ||   | ||  |  |_|  ||    _  ||_     _|
|       ||   |_||_ |       ||   |_| |  |   |
|      _||    __  ||_     _||    ___|  |   |
|     |_ |   |  | |  |   |  |   |      |   |
|_______||___|  |_|  |___|  |___|      |___|`

var Command = &cobra.Command{
	Use:   "info",
	Short: "information about your crypt file",
	Long:  `displays meta information about your crypt file.`,
	RunE:  info,
}

func info(cmd *cobra.Command, args []string) error {
	env, err := environment.Load(cmd)
	if err != nil {
		return err
	}
	defer env.Close()

	crypt := env.Crypt()

	data := [][]string{
		{"id", crypt.ID},
		{"credentials", strconv.Itoa(crypt.TotalActiveCredentials)},
		{"created at", crypt.CreatedAt.Format("01/02/2006")},
		{"updated at", crypt.UpdatedAt.Format("01/02/2006")},
	}
	fmt.Println(title)
	utils.PrintTable(data, utils.TableOpts{})

	return nil
}
