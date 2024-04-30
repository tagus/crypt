package environment

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/ciphers/aescipher"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/crypt/internal/repos/dbrepo"
	"github.com/tagus/mango"
	"os"
	"path/filepath"
)

type Environment struct {
	repo  repos.Repo
	crypt *repos.Crypt
}

type InitStoreOpts struct {
	CryptName   string
	CryptDBPath string
	InitFile    bool
}

func initEnv(ctx context.Context, opts InitStoreOpts) (*Environment, error) {
	mango.Debug("initializing environment:", opts)

	path, err := resolveCryptDBPath(opts.CryptDBPath, !opts.InitFile)
	if err != nil {
		return nil, err
	}
	mango.Debug("using crypt at: ", color.YellowString(path))

	ci, err := initCipher()
	if err != nil {
		return nil, err
	}

	repo, err := dbrepo.Initialize(ctx, path, ci)
	if err != nil {
		return nil, err
	}

	st := &Environment{repo: repo}

	crypts, err := repo.QueryCrypts(ctx, repos.QueryCryptsFilter{Name: opts.CryptName})
	if err != nil {
		return nil, err
	}
	if len(crypts) != 1 {
		crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{ID: mango.ShortID(), Name: opts.CryptName})
		if err != nil {
			return nil, err
		}
		st.crypt = crypt
	} else {
		st.crypt = crypts[0]
	}

	return st, nil
}

func initCipher() (ciphers.Cipher, error) {
	ak := asker.DefaultAsker()
	pwd, err := ak.AskSecret(color.YellowString("pwd"), false)
	if err != nil {
		return nil, err
	}
	return aescipher.New(pwd)
}

// resolveCryptDBPath determines the path of the crypt db to be used, the `crypt-db`
// flag takes priority, falling back to a `CRYPT_DB` env var, `.crypt.db` in the current working
// directory and finally defaulting to a `.crypt.db` in the current user's home directory
func resolveCryptDBPath(cryptDBPath string, checkPath bool) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	hd, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	paths := []string{
		cryptDBPath,
		os.Getenv("CRYPT_DB"),
		filepath.Join(wd, ".crypt.db"),
		filepath.Join(hd, ".crypt.db"),
	}

	for _, path := range paths {
		if path == "" {
			continue
		}
		if checkPath && !mango.FileExists(path) {
			continue
		}
		return path, nil
	}

	return "", errors.New("no valid crypt db file found")
}

func (e *Environment) Crypt() *repos.Crypt {
	return e.crypt
}

func (e *Environment) Repo() repos.Repo {
	return e.repo
}
