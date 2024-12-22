package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/cli/edit"

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

const (
	Version     = "v2.1.2"
	VerboseFlag = "verbose"
	LogPrefix   = "crypt"
)

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
	SilenceUsage:      true,
	SilenceErrors:     true,
	Version:           Version,
	PersistentPreRunE: initialize,
}

func Execute() {

	rootCmd.
		PersistentFlags().
		StringP(environment.CryptDBPathFlag, "c", "", "the crypt db location")
	rootCmd.
		PersistentFlags().
		Bool(environment.CryptDBInitFlag, false, "whether to initialize the crypt db file if it doesn't exist")
	rootCmd.
		PersistentFlags().
		StringP(environment.CryptNameFlag, "n", "main", "the crypt name")
	rootCmd.
		PersistentFlags().
		BoolP(VerboseFlag, "v", false, "print out any debug logs")

	// adding all subcommands
	rootCmd.AddCommand(info.Command)
	rootCmd.AddCommand(add.Command)
	rootCmd.AddCommand(show.Command)
	rootCmd.AddCommand(archive.Command)
	rootCmd.AddCommand(list.Command)
	rootCmd.AddCommand(edit.Command)

	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, ciphers.ErrInvalidPassword) {
			fmt.Println(color.YellowString("invalid password"))
		} else if errors.Is(err, asker.ErrInterrupt) {
			os.Exit(0)
		} else {
			mango.Fatal(err)
		}
	}
}

func initialize(cmd *cobra.Command, args []string) error {
	// checking log level
	isVerbose, err := cmd.Flags().GetBool(VerboseFlag)
	if err != nil {
		return err
	}

	if isVerbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("log level is set to verbose")
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	return nil
}
