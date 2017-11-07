package cmd

import (
	"fmt"
	// "os"
	// "io/ioutil"
	"github.com/spf13/cobra"
	// "golang.org/x/crypto/ssh/terminal"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a credential to crypt",
	Long: `crypt add [service] [username] [password]`,
	Run: add,
}

func init() {
	RootCmd.AddCommand(addCmd)
}

func add(cmd *cobra.Command, args []string) {
	fmt.Println("")
	// _, err := terminal.ReadPassword(0)
	// if err != nil {
	// 	panic(err.Error())
	// 	os.Exit(1)
	// } else {
	// 	// use pwd to decrypt JSON file
	// 	cipherText, _ := ioutil.ReadFile(cryptFile)
	// 	fmt.Println(cipherText)
	// }
}
