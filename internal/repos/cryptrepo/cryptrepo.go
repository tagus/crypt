package cryptrepo

import (
	"context"

	"github.com/tagus/crypt/internal/repos/dbrepo"

	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/repos"
)

// CryptRepo scopes the underlying repo to the given crypt id and cipher
type CryptRepo struct {
	repo    *dbrepo.DbRepo
	cryptID string
	ci      ciphers.Cipher
}

func New(repo *dbrepo.DbRepo, cryptID string, ci ciphers.Cipher) *CryptRepo {
	return &CryptRepo{
		repo:    repo,
		cryptID: cryptID,
		ci:      ci,
	}
}

func (c *CryptRepo) QueryCredentials(
	ctx context.Context,
	filter repos.QueryCredentialsFilter,
) ([]*repos.Credential, error) {
	filter.CryptID = c.cryptID
	return c.repo.QueryCredentials(ctx, c.ci, filter)
}

func (c *CryptRepo) InsertCredential(
	ctx context.Context,
	cred *repos.Credential,
) (*repos.Credential, error) {
	return c.repo.InsertCredential(ctx, c.ci, c.cryptID, cred)
}

func (c *CryptRepo) UpdateCredential(
	ctx context.Context,
	cred *repos.Credential,
) (*repos.Credential, error) {
	return c.repo.UpdateCredential(ctx, c.ci, c.cryptID, cred)
}

func (c *CryptRepo) AccessCredential(
	ctx context.Context,
	credID string,
) (*repos.Credential, error) {
	return c.repo.AccessCredential(ctx, c.ci, c.cryptID, credID)
}

func (c *CryptRepo) ArchiveCredential(
	ctx context.Context,
	credID string,
) error {
	return c.repo.ArchiveCredential(ctx, credID)
}
