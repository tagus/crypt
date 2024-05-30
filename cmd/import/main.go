package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/asker"
	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/ciphers/aescipher"
	"github.com/tagus/crypt/internal/crypt"
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
	mango.Init(mango.LogLevelDebug, "import")

	if *cp == "" {
		mango.Fatal("crypt file path is required")
	}
	mango.Debug("crypt file path: ", *cp)
	if *db == "" {
		mango.Fatal("sqlite db file path is required")
	}
	mango.Debug("sqlite db file path: ", *db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo, err := dbrepo.Initialize(ctx, *db)
	mango.FatalIf(err)

	/******************************************************************************/

	// asking user for master password that will be used to decrypt the secure values
	ak := asker.DefaultAsker()
	pwd, err := ak.AskSecret(color.YellowString("pwd"), true)
	mango.FatalIf(err)

	signature := []byte(ciphers.ComputeHash(mango.ShortID()))
	hashedPwd, err := ciphers.ComputeHashPwd(pwd)
	mango.FatalIf(err)

	ci, err := aescipher.New(pwd, hashedPwd, signature)
	mango.FatalIf(err)

	/******************************************************************************/

	cr, err := parseCryptFile(*cp)
	mango.FatalIf(err)

	mango.Debug("importing crypt:", cr.Id)
	newCrypt, err := repo.InsertCrypt(
		ctx,
		&repos.Crypt{ID: cr.Id, Name: *name, HashedPassword: hashedPwd, Signature: signature},
	)
	mango.FatalIf(err)

	/******************************************************************************/

	for _, cred := range cr.Credentials {
		newCred, err := repo.InsertCredential(ctx, ci, newCrypt.ID, &repos.Credential{
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

	/******************************************************************************/

	mango.Debug("import complete for crypt:", newCrypt.ID)
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
