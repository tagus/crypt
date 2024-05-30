package cli

import (
	"errors"
	"os"

	"github.com/tagus/crypt/internal/cli/list"

	"github.com/tagus/crypt/internal/cli/archive"

	"github.com/spf13/cobra"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/cli/add"
	"github.com/tagus/crypt/internal/cli/environment"
	"github.com/tagus/crypt/internal/cli/info"
	"github.com/tagus/crypt/internal/cli/show"
	"github.com/tagus/mango"
)

const Version = "v2.0.0"

var rootCmd = &cobra.Command{
	Use:   "crypt",
	Short: "a secure credential store",
	Long: `crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts.

Crypt stores your credentials in a sqlite db that will encrypt private details
(e.g. password, secret keys, security questions) with a master password so that 
they cannot be read even if the sqlite file is examined.

The db file can be specified using the following methods listed here in decreasing priority.

	1. db flag
	2. CRYPT_DB env variable
	3. ./.crypt.db
	4. ~/.crypt.db`,
	SilenceUsage: true,
	Version:      Version,
}

func Execute() {
	mango.Init(mango.LogLevelDebug, "crypt")

	rootCmd.
		PersistentFlags().
		StringP(environment.CryptDBPathFlag, "c", "", "the crypt db location")
	rootCmd.
		PersistentFlags().
		Bool(environment.CryptDBInitFlag, false, "whether to initialize the crypt db file if it doesn't exist")
	rootCmd.
		PersistentFlags().
		StringP(environment.CryptNameFlag, "n", "main", "the crypt name")

	// adding all subcommands
	rootCmd.AddCommand(info.Command)
	rootCmd.AddCommand(add.Command)
	rootCmd.AddCommand(show.Command)
	rootCmd.AddCommand(archive.Command)
	rootCmd.AddCommand(list.Command)

	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, asker.ErrInterrupt) {
			os.Exit(0)
		} else {
			mango.Fatal(err)
		}
	}
}
