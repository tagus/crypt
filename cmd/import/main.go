package main

import (
	"context"
	"flag"
	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/ciphers/aescipher"
	"github.com/tagus/crypt/internal/crypt"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/crypt/internal/repos/dbrepo"
	"github.com/tagus/mango"
	"os"
)

var (
	cp   = flag.String("crypt", "", "crypt file path")
	db   = flag.String("db", "", "sqlite db file path")
	name = flag.String("name", "main", "crypt name")
)

func main() {
	flag.Parse()
	mango.Init(mango.LogLevelDebug, "import")

	if *cp == "" {
		mango.Fatal("crypt file path is required")
	}
	mango.Debug("crypt file path: ", *cp)
	if *db == "" {
		mango.Fatal("sqlite db file path is required")
	}
	mango.Debug("sqlite db file path: ", *db)

	ctx := context.Background()

	ci, err := initCipher()
	mango.FatalIf(err)

	repo, err := dbrepo.Initialize(ctx, *db, ci)
	mango.FatalIf(err)

	cr, err := parseCryptFile(*cp)
	mango.FatalIf(err)

	mango.Debug("importing crypt:", cr.Id)
	newCrypt, err := getOrCreateCrypt(ctx, repo, cr, *name)
	mango.FatalIf(err)

	for _, cred := range cr.Credentials {
		newCred, err := repo.InsertCredential(ctx, newCrypt.ID, &repos.Credential{
			ID:          cred.Id,
			Service:     cred.Service,
			Email:       cred.Email,
			Username:    cred.Username,
			Password:    cred.Password,
			Description: cred.Description,
			Details:     &repos.Details{},
			Tags:        cred.Tags,
			Domains:     []string{},
			CreatedAt:   cred.GetCreatedAt(),
			UpdatedAt:   cred.GetUpdatedAt(),
		})
		mango.FatalIf(err)
		mango.Debug("imported service:", newCred.Service)
	}

	mango.Debug("import complete for crypt:", newCrypt.ID)
}

func getOrCreateCrypt(ctx context.Context, repo repos.Repo, cr *crypt.Crypt, name string) (*repos.Crypt, error) {
	crypts, err := repo.QueryCrypts(ctx, repos.QueryCryptsFilter{ID: cr.Id})
	if err != nil {
		return nil, err
	}
	if len(crypts) != 1 {
		return repo.InsertCrypt(ctx, &repos.Crypt{ID: cr.Id, Name: name})
	}
	return crypts[0], nil
}

func initCipher() (ciphers.Cipher, error) {
	ak := asker.DefaultAsker()
	pwd, err := ak.AskSecret(color.YellowString("pwd"), true)
	if err != nil {
		return nil, err
	}
	return aescipher.New(pwd)
}

func parseCryptFile(path string) (*crypt.Crypt, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cr, err := crypt.FromJSON(buf)
	if err != nil {
		return nil, err
	}
	return cr, nil
}
