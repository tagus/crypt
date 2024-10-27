package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/ciphers/aescipher"
	"github.com/tagus/crypt/internal/legacycrypt"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/crypt/internal/repos/dbrepo"
	"github.com/tagus/mango"
)

var (
	cp   = flag.String("crypt", "", "crypt file path")
	db   = flag.String("db", "", "sqlite db file path")
	name = flag.String("name", "main", "crypt name")
)

func main() {
	flag.Parse()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if *cp == "" {
		mango.Fatal("crypt file path is required")
	}
	slog.Debug("crypt file path ", "path", *cp)
	if *db == "" {
		mango.Fatal("sqlite db file path is required")
	}
	slog.Debug("sqlite db file path", "path", *db)

	/******************************************************************************/

	// asking user for master password that will be used to decrypt the secure values
	ak := asker.DefaultAsker()
	pwd, err := ak.AskSecret(color.YellowString("pwd"), true)
	mango.FatalIf(err)
	slog.Debug("collected password")

	signature := []byte(ciphers.ComputeHash(mango.ShortID()))
	hashedPwd, err := ciphers.ComputeHashPwd(pwd)
	mango.FatalIf(err)

	ci, err := aescipher.New(pwd, hashedPwd, signature)
	mango.FatalIf(err)

	/******************************************************************************/

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	repo, err := dbrepo.Initialize(ctx, *db)
	mango.FatalIf(err)

	cr, err := parseCryptFile(*cp)
	mango.FatalIf(err)

	slog.Debug("importing crypt", "id", cr.Id)
	newCrypt, err := repo.InsertCrypt(
		ctx,
		&repos.Crypt{
			ID:             cr.Id,
			Name:           *name,
			HashedPassword: hashedPwd,
			Signature:      signature,
			CreatedAt:      cr.GetCreatedAt(),
			UpdatedAt:      cr.GetUpdatedAt(),
		},
	)
	mango.FatalIf(err)

	/******************************************************************************/

	for _, cred := range cr.Credentials {
		newCred, err := repo.InsertCredential(ctx, ci, newCrypt.ID, &repos.Credential{
			ID:            cred.Id,
			Service:       cred.Service,
			Email:         cred.Email,
			Username:      cred.Username,
			Password:      cred.Password,
			Description:   cred.Description,
			Details:       &repos.Details{},
			Tags:          cred.Tags,
			Domains:       []string{},
			CreatedAt:     cred.GetCreatedAt(),
			UpdatedAt:     cred.GetUpdatedAt(),
			AccessedAt:    cred.GetAccessedAt(),
			AccessedCount: cred.AccessedCount,
		})
		mango.FatalIf(err)
		slog.Debug("imported service", "service", newCred.Service, "accessed_count", cred.AccessedCount)
	}

	/******************************************************************************/

	slog.Debug("import complete for crypt", "id", newCrypt.ID)
}

func parseCryptFile(path string) (*legacycrypt.Crypt, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cr, err := legacycrypt.FromJSON(buf)
	if err != nil {
		return nil, err
	}
	return cr, nil
}
