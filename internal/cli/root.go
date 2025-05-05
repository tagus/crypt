package cli

import (
	"errors"
	"log/slog"
	"os"

	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/cli/cutils"
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
	Version     = "v2.1.4"
	VerboseFlag = "verbose"
	AppLabel    = "crypt"
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
	4. ~/.crypt.db,
	5. ~/.config/crypt.db`,
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
			slog.Warn("invalid password")
		} else if errors.Is(err, cutils.ErrNoCredentialFound) {
			slog.Warn("invalid service provided")
		} else if errors.Is(err, asker.ErrInterrupt) {
			os.Exit(0)
		} else {
			isVerbose, _ := rootCmd.Flags().GetBool(VerboseFlag)
			if isVerbose {
				slog.Error(err.Error())
			} else {
				mango.Fatal(err)
			}
		}
	}
}

func initialize(cmd *cobra.Command, args []string) error {
	isVerbose, err := cmd.Flags().GetBool(VerboseFlag)
	if err != nil {
		return err
	}
	level := slog.LevelInfo
	if isVerbose {
		level = slog.LevelDebug
	}
	mango.Init(level, AppLabel)
	return nil
}
