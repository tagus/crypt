package environment

import (
	"github.com/spf13/cobra"
)

const (
	CryptDBPathFlag = "crypt-db"
	CryptDBInitFlag = "init"
	CryptNameFlag   = "crypt-name"
)

var (
	store *Environment
)

func Load(cmd *cobra.Command) (*Environment, error) {
	if store != nil {
		return store, nil
	}

	cryptDB, err := cmd.Flags().GetString(CryptDBPathFlag)
	if err != nil {
		return nil, err
	}

	initCryptDBFile, err := cmd.Flags().GetBool(CryptDBInitFlag)
	if err != nil {
		return nil, err
	}

	cryptName, err := cmd.Flags().GetString(CryptNameFlag)
	if err != nil {
		return nil, err
	}

	st, err := initEnv(cmd.Context(), InitStoreOpts{
		CryptName:   cryptName,
		CryptDBPath: cryptDB,
		InitFile:    initCryptDBFile,
	})
	if err != nil {
		return nil, err
	}
	store = st
	return store, nil
}
