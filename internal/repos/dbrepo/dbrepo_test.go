package dbrepo

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/tagus/crypt/internal/ciphers/passthroughcipher"

	"github.com/stretchr/testify/require"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/mango"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	repo *DbRepo
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("sqlite3", "./dbrepo_test.db")
	mango.FatalIf(err)

	ctx := context.TODO()
	repo, err = initialize(ctx, db)
	mango.FatalIf(err)

	resetDB()
	exitCode := m.Run()
	mango.FatalIf(os.Remove("./dbrepo_test.db"))
	os.Exit(exitCode)
}

func resetDB() {
	_, err := db.Exec("DELETE FROM crypts")
	mango.FatalIf(err)
	_, err = db.Exec("DELETE FROM credentials")
	mango.FatalIf(err)
	_, err = db.Exec("DELETE FROM credential_versions")
	mango.FatalIf(err)
}

func TestDbRepo_QueryCrypts(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	_, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	_, err = repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-2",
		Name:           "alt_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	crypts, err := repo.QueryCrypts(ctx, repos.QueryCryptsFilter{})
	require.NoError(t, err)
	require.Len(t, crypts, 2)

	crypts, err = repo.QueryCrypts(ctx, repos.QueryCryptsFilter{
		ID: "test-crypt-1",
	})
	require.NoError(t, err)
	require.Len(t, crypts, 1)
	require.Equal(t, "test-crypt-1", crypts[0].ID)
	require.Equal(t, "default_crypt", crypts[0].Name)
	require.Equal(t, []byte("hashed_password"), crypts[0].HashedPassword)
	require.Equal(t, []byte("signature"), crypts[0].Signature)

	crypts, err = repo.QueryCrypts(ctx, repos.QueryCryptsFilter{
		Name: "alt_crypt",
	})
	require.NoError(t, err)
	require.Len(t, crypts, 1)

	crypts, err = repo.QueryCrypts(ctx, repos.QueryCryptsFilter{
		Name: "alt_crypt_2",
	})
	require.NoError(t, err)
	require.Len(t, crypts, 0)
}

func TestDbRepo_CreateCrypt(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	require.Equal(t, "test-crypt-1", crypt.ID)
	require.Equal(t, "default_crypt", crypt.Name)
	require.NotEmpty(t, crypt.CreatedAt)
	require.NotEmpty(t, crypt.UpdatedAt)
}

func TestDbRepo_InsertCredential(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	ci := passthroughcipher.New()
	cred, err := repo.InsertCredential(ctx, ci, crypt.ID, &repos.Credential{
		ID:          "credential-1",
		Service:     "test-service",
		Domains:     []string{"domain-1", "domain-2"},
		Email:       "test@test.com",
		Username:    "username",
		Password:    "password",
		Description: "description",
		Details: &repos.Details{
			SecurityQuestions: []repos.SecurityQuestion{
				{
					Question: "question",
					Answer:   "answer",
				},
			},
		},
		Tags: []string{"tag-1", "tag-2"},
	})
	require.NoError(t, err)
	require.Equal(t, "credential-1", cred.ID)
	require.Equal(t, "test-service", cred.Service)
	require.Equal(t, []string{"domain-1", "domain-2"}, cred.Domains)
	require.Equal(t, "test@test.com", cred.Email)
	require.Equal(t, "username", cred.Username)
	require.Equal(t, "password", cred.Password)
	require.Equal(t, "description", cred.Description)
	require.Equal(t, []string{"tag-1", "tag-2"}, cred.Tags)
	require.Equal(t, 1, cred.Version)

	require.Len(t, cred.Details.SecurityQuestions, 1)
	require.Equal(t, cred.Details.SecurityQuestions[0].Question, "question")
	require.Equal(t, cred.Details.SecurityQuestions[0].Answer, "answer")
}

func TestDbRepo_InsertCredentialWithoutID(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	ci := passthroughcipher.New()
	_, err = repo.InsertCredential(ctx, ci, crypt.ID, &repos.Credential{
		Service:     "test-service",
		Domains:     nil,
		Email:       "test@test.com",
		Username:    "username",
		Password:    "password",
		Description: "description",
		Details:     nil,
		Tags:        nil,
	})
	require.Error(t, err)
}

func TestDbRepo_InsertCredentialWithNullValues(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	ci := passthroughcipher.New()
	cred, err := repo.InsertCredential(ctx, ci, crypt.ID, &repos.Credential{
		ID:          "credential-1",
		Service:     "test-service",
		Domains:     nil,
		Email:       "test@test.com",
		Username:    "username",
		Password:    "password",
		Description: "description",
		Details:     nil,
		Tags:        nil,
	})
	require.NoError(t, err)
	require.NotNil(t, cred)

	var details *repos.Details
	require.Len(t, cred.Domains, 0)
	require.Equal(t, details, cred.Details)
}

func TestDbRepo_UpdateCredential(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	ci := passthroughcipher.New()
	cred, err := repo.InsertCredential(ctx, ci, crypt.ID, &repos.Credential{
		ID:          "credential-1",
		Service:     "test-service",
		Domains:     []string{"domain-1", "domain-2"},
		Email:       "test@test.com",
		Username:    "username",
		Password:    "password",
		Description: "description",
		Details: &repos.Details{
			SecurityQuestions: []repos.SecurityQuestion{
				{
					Question: "question",
					Answer:   "answer",
				},
			},
		},
		Tags: []string{"tag-1", "tag-2"},
	})
	require.NoError(t, err)
	require.Equal(t, "credential-1", cred.ID)
	require.Equal(t, "test-service", cred.Service)
	require.Equal(t, "password", cred.Password)
	require.Equal(t, 1, cred.Version)

	cred.Password = "new-password"
	cred, err = repo.UpdateCredential(ctx, ci, crypt.ID, cred)
	require.NoError(t, err)
	require.Equal(t, "new-password", cred.Password)
	require.Equal(t, 2, cred.Version)
}

func TestDbRepo_QueryCredentialsWithIncrementAccessCountFlag(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	ci := passthroughcipher.New()
	cred, err := repo.InsertCredential(ctx, ci, crypt.ID, &repos.Credential{
		ID:          "credential-1",
		Service:     "test-service",
		Domains:     []string{"domain-1", "domain-2"},
		Email:       "test@test.com",
		Username:    "username",
		Password:    "password",
		Description: "description",
		Details: &repos.Details{
			SecurityQuestions: []repos.SecurityQuestion{
				{
					Question: "question",
					Answer:   "answer",
				},
			},
		},
		Tags: []string{"tag-1", "tag-2"},
	})
	require.NoError(t, err)
	require.Equal(t, 0, cred.AccessedCount)

	creds, err := repo.QueryCredentials(
		ctx,
		ci,
		repos.QueryCredentialsFilter{ID: cred.ID, IncrementAccessCount: true},
	)
	require.NoError(t, err)
	require.Len(t, creds, 1)
	require.Equal(t, 1, creds[0].AccessedCount)
	require.NotNil(t, creds[0].AccessedAt)

	newCreds, err := repo.QueryCredentials(ctx, ci, repos.QueryCredentialsFilter{ID: cred.ID})
	require.NoError(t, err)
	require.Len(t, newCreds, 1)
	require.Equal(t, creds[0].AccessedCount, newCreds[0].AccessedCount)
	require.NotNil(t, newCreds[0].AccessedAt)
}

func TestDbRepo_ArchiveCredential(t *testing.T) {
	defer resetDB()

	ctx := context.TODO()
	crypt, err := repo.InsertCrypt(ctx, &repos.Crypt{
		ID:             "test-crypt-1",
		Name:           "default_crypt",
		HashedPassword: []byte("hashed_password"),
		Signature:      []byte("signature"),
	})
	require.NoError(t, err)

	ci := passthroughcipher.New()
	cred, err := repo.InsertCredential(ctx, ci, crypt.ID, &repos.Credential{
		ID:          "credential-1",
		Service:     "test-service",
		Domains:     []string{"domain-1", "domain-2"},
		Email:       "test@test.com",
		Username:    "username",
		Password:    "password",
		Description: "description",
		Details: &repos.Details{
			SecurityQuestions: []repos.SecurityQuestion{
				{
					Question: "question",
					Answer:   "answer",
				},
			},
		},
		Tags: []string{"tag-1", "tag-2"},
	})
	require.NoError(t, err)
	require.Equal(t, 0, cred.AccessedCount)

	err = repo.ArchiveCredential(ctx, cred.ID)
	require.NoError(t, err)

	creds, err := repo.QueryCredentials(ctx, ci, repos.QueryCredentialsFilter{ID: cred.ID})
	require.NoError(t, err)
	require.Len(t, creds, 0)
}
