package cmd

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

var key, cryptFile string

var RootCmd = &cobra.Command{
	Use:   "crypt",
	Short: "A secure credential store",
	Long:`Crypt is CLI application to securely store your credentials
	so that you don't have to worry about remembering all of your
	internet accounts`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initCrypt)
	// RootCmd.PersistentFlags().StringVar(&cryptFile, "crypt", "", "cipher text credential store (default is $HOME/.crypt)")
}

type Credential struct {
	Username			string
	Password			string
}

type Crypt struct {
	Credentials []Credential
}

func initCrypt() {
	home, _ := homedir.Dir()
	cryptFilePath := fmt.Sprintf("%s/.crypt", home)

	if _, err := os.Stat(cryptFilePath); os.IsNotExist(err) {
		color.Yellow("Crypt file not found")
		fmt.Print("Would you like to create one [y/n] ")
		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "yes" {
			// cryptFile, err := os.Create(cryptFilePath)
			// if err != nil {
			// 	color.Red("There was an error during crypt file creation")
			// 	os.Exit(1)
			// }

			color.Green("Crypt file created at '%s'", cryptFilePath)
			credentials := []Credential {
				Credential {"sdk282", "s@f3pa$$w0rd"},
				Credential {"tom233", "apple"},
			}
			crypt := Crypt{credentials}
			cryptData, _ := json.Marshal(crypt)
			// cryptData, _ := json.Marshal(Credential{"tom@gmail.com", "burrito"})

			fmt.Printf("%s\n", cryptData)
			var cryptDecoded Crypt
			if err := json.Unmarshal(cryptData, &cryptDecoded); err != nil {
				color.Red("There was an error unmarshalling data: %s", err)
				os.Exit(1)
			}

			fmt.Println(cryptDecoded)
			// cryptFile.Close()
		} else {
			os.Exit(0)
		}
	} else {
		color.Green("Crypt file found")
	}
}