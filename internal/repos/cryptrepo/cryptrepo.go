package cryptrepo

import (
	"context"

	"github.com/tagus/crypt/internal/ciphers"
	"github.com/tagus/crypt/internal/repos"
)

type Repo interface {
	QueryCredentials(
		ctx context.Context,
		ci ciphers.Cipher,
		filter repos.QueryCredentialsFilter,
	) ([]*repos.Credential, error)
	InsertCredential(
		ctx context.Context,
		ci ciphers.Cipher,
		cryptID string,
		cred *repos.Credential,
	) (*repos.Credential, error)
	UpdateCredential(
		ctx context.Context,
		ci ciphers.Cipher,
		cryptID string,
		cred *repos.Credential,
	) (*repos.Credential, error)
	AccessCredential(
		ctx context.Context,
		ci ciphers.Cipher,
		cryptID, credID string,
	) (*repos.Credential, error)
}

// CryptRepo scopes the underlying repo to the given crypt id and cipher
type CryptRepo struct {
	repo    Repo
	cryptID string
	ci      ciphers.Cipher
}

func New(repo Repo, cryptID string, ci ciphers.Cipher) *CryptRepo {
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
