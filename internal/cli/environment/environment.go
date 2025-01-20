package environment

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/tagus/crypt/internal/repos/cryptrepo"

	"github.com/mitchellh/go-homedir"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/ciphers/aescipher"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/crypt/internal/repos/dbrepo"
	"github.com/tagus/mango"
)

type Environment struct {
	repo  *dbrepo.DbRepo
	cr    CryptRepo
	crypt *repos.Crypt
}

type CryptRepo interface {
	QueryCredentials(
		ctx context.Context,
		filter repos.QueryCredentialsFilter,
	) ([]*repos.Credential, error)
	InsertCredential(ctx context.Context, cred *repos.Credential) (*repos.Credential, error)
	UpdateCredential(ctx context.Context, cred *repos.Credential) (*repos.Credential, error)
	AccessCredential(ctx context.Context, credID string) (*repos.Credential, error)
	ArchiveCredential(ctx context.Context, credID string) error
}

type InitStoreOpts struct {
	CryptName   string
	CryptDBPath string
	InitFile    bool
}

func initEnv(ctx context.Context, opts InitStoreOpts) (*Environment, error) {
	slog.Debug("initializing environment", "opts", opts)

	path, err := resolveCryptDBPath(opts.CryptDBPath, !opts.InitFile)
	if err != nil {
		return nil, err
	}
	slog.Debug("using crypt", "path", path, "name", opts.CryptName)

	repo, err := dbrepo.Initialize(ctx, path)
	if err != nil {
		return nil, err
	}

	/******************************************************************************/

	ak := asker.DefaultAsker()
	pwd, err := ak.AskSecret(mango.ColorizeYellow("pwd"), false)
	if err != nil {
		return nil, err
	}

	var crypt *repos.Crypt
	crypts, err := repo.QueryCrypts(ctx, repos.QueryCryptsFilter{Name: opts.CryptName})
	if err != nil {
		return nil, err
	}
	if len(crypts) != 1 {
		signature := []byte(ciphers.ComputeHash(mango.ShortID()))
		hashedPwd, err := ciphers.ComputeHashPwd(pwd)
		if err != nil {
			return nil, err
		}
		_crypt, err := repo.InsertCrypt(
			ctx,
			&repos.Crypt{
				ID:             mango.ShortID(),
				Name:           opts.CryptName,
				HashedPassword: hashedPwd,
				Signature:      signature,
			},
		)
		if err != nil {
			return nil, err
		}
		crypt = _crypt
	} else {
		crypt = crypts[0]
	}

	/******************************************************************************/

	ci, err := aescipher.New(pwd, crypt.HashedPassword, crypt.Signature)
	if err != nil {
		return nil, err
	}

	return &Environment{
		repo:  repo,
		cr:    cryptrepo.New(repo, crypt.ID, ci),
		crypt: crypt,
	}, nil
}

// resolveCryptDBPath determines the path of the crypt db to be used in the following order:
//   - the `crypt-db` flag
//   - `CRYPT_DB` env var
//   - ./.crypt.db
//   - ~/.crypt.db
//   - ~/.config/crypt/crypt.db
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
		filepath.Join(hd, ".config", "crypt", "crypt.db"),
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

func (e *Environment) Repo() CryptRepo {
	return e.cr
}

func (e *Environment) Close() error {
	return e.repo.Close()
}
