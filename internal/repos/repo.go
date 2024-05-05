package repos

import "context"

type QueryCryptsFilter struct {
	Name string
	ID   string
}

type QueryCredentialsFilter struct {
	ID      string
	CryptID string
	Service string
}

type Repo interface {
	QueryCrypts(ctx context.Context, filter QueryCryptsFilter) ([]*Crypt, error)
	InsertCrypt(ctx context.Context, crypt *Crypt) (*Crypt, error)
	QueryCredentials(ctx context.Context, filter QueryCredentialsFilter) ([]*Credential, error)
	InsertCredential(ctx context.Context, cryptID string, cred *Credential) (*Credential, error)
	UpdateCredential(ctx context.Context, cryptID string, cred *Credential) (*Credential, error)
	Close() error
}
