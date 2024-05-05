package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/tagus/crypt/internal/ciphers"
	"time"

	_ "embed"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/mango"
)

type DbRepo struct {
	cipher ciphers.Cipher
	db     *sql.DB
}

//go:embed schema.sql
var schema string

// Initialize will ensure that the corresponding sqlite file contains the crypt tables
func Initialize(ctx context.Context, path string, cipher ciphers.Cipher) (*DbRepo, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return initialize(ctx, cipher, db)
}

func initialize(ctx context.Context, cipher ciphers.Cipher, db *sql.DB) (*DbRepo, error) {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, err
	}
	return &DbRepo{
		db:     db,
		cipher: cipher,
	}, nil
}

func (r *DbRepo) Close() error {
	return r.db.Close()
}

/******************************************************************************/

func (r *DbRepo) QueryCrypts(ctx context.Context, filter repos.QueryCryptsFilter) ([]*repos.Crypt, error) {
	qb := sq.Select(
		"id",
		"name",
		"updated_at",
		"created_at",
		"total_active_credentials",
	).
		From("crypts").
		Where(sq.Eq{"archived_at": nil}).
		RunWith(r.db)

	if filter.Name != "" {
		qb = qb.Where(sq.Eq{"name": filter.Name})
	}
	if filter.ID != "" {
		qb = qb.Where(sq.Eq{"id": filter.ID})
	}

	rows, err := qb.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var crypts []*repos.Crypt
	for rows.Next() {
		crypt, err := r.parseCrypt(rows)
		if err != nil {
			return nil, err
		}
		crypts = append(crypts, crypt)
	}

	return crypts, nil
}

func (r *DbRepo) parseCrypt(row *sql.Rows) (*repos.Crypt, error) {
	var crypt repos.Crypt

	err := row.Scan(&crypt.ID, &crypt.Name, &crypt.UpdatedAt, &crypt.CreatedAt, &crypt.TotalActiveCredentials)
	if err != nil {
		return nil, err
	}
	return &crypt, nil
}

/******************************************************************************/

func (r *DbRepo) InsertCrypt(ctx context.Context, crypt *repos.Crypt) (*repos.Crypt, error) {
	qb := sq.Insert("crypts").
		Columns("id", "name").
		Values(crypt.ID, crypt.Name).
		RunWith(r.db)

	if _, err := qb.ExecContext(ctx); err != nil {
		return nil, err
	}

	crypts, err := r.QueryCrypts(ctx, repos.QueryCryptsFilter{
		ID: crypt.ID,
	})
	if err != nil {
		return nil, err
	}
	if len(crypts) == 0 {
		return nil, errors.New("failed to create the given crypt")
	}

	return crypts[0], nil
}

/******************************************************************************/

func (r *DbRepo) QueryCredentials(ctx context.Context, filter repos.QueryCredentialsFilter) ([]*repos.Credential, error) {
	qb := sq.Select(
		"cr.id",
		"cr.tags",
		"cr.updated_at",
		"cr.created_at",
		"cr.accessed_at",
		"cr.accessed_count",
		"cv.service",
		"cv.email",
		"cv.domains",
		"cv.username",
		"cv.s_password",
		"cv.description",
		"cv.s_details",
		"cv.version",
	).
		From("credentials AS cr").
		InnerJoin("credential_versions AS cv ON cr.id = cv.credential_id AND cr.current_version = cv.version").
		Where(sq.Eq{"cr.archived_at": nil}).
		RunWith(r.db)

	if filter.CryptID != "" {
		qb = qb.Where(sq.Eq{"cr.crypt_id": filter.CryptID})
	}
	if filter.Service != "" {
		qb = qb.Where("cv.service LIKE ?", fmt.Sprintf("%%%s%%", filter.Service))
	}
	if filter.ID != "" {
		qb = qb.Where(sq.Eq{"cr.id": filter.ID})
	}

	rows, err := qb.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var creds []*repos.Credential
	for rows.Next() {
		cred, err := r.parseCredential(rows)
		if err != nil {
			return nil, err
		}
		creds = append(creds, cred)
	}

	return creds, nil
}

func (r *DbRepo) parseCredential(row *sql.Rows) (*repos.Credential, error) {
	var (
		cred             repos.Credential
		tagsJSON         string
		accessedAt       sql.NullTime
		domainsJSON      string
		encryptedPwd     []byte
		encryptedDetails []byte
	)

	err := row.Scan(
		&cred.ID,
		&tagsJSON,
		&cred.UpdatedAt,
		&cred.CreatedAt,
		&accessedAt,
		&cred.AccessedCount,
		&cred.Service,
		&cred.Email,
		&domainsJSON,
		&cred.Username,
		&encryptedPwd,
		&cred.Description,
		&encryptedDetails,
		&cred.Version,
	)
	if err != nil {
		return nil, err
	}

	domains, err := mango.UnmarshalFromString[[]string](domainsJSON)
	if err != nil {
		return nil, err
	}
	cred.Domains = *domains

	tags, err := mango.UnmarshalFromString[[]string](tagsJSON)
	if err != nil {
		return nil, err
	}
	cred.Tags = *tags

	cred.Password, err = r.cipher.Decrypt(encryptedPwd)
	if err != nil {
		return nil, err
	}

	detailsJSON, err := r.cipher.Decrypt(encryptedDetails)
	if err != nil {
		return nil, err
	}
	details, err := mango.UnmarshalFromString[repos.Details](detailsJSON)
	if err != nil {
		return nil, err
	}
	cred.Details = details

	return &cred, nil
}

/******************************************************************************/

func (r *DbRepo) InsertCredential(ctx context.Context, cryptID string, cred *repos.Credential) (*repos.Credential, error) {
	tagsJSON, err := mango.MarshalToString(cred.Tags)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	_, err = sq.Insert("credentials").
		Columns(
			"id",
			"current_version",
			"latest_version",
			"tags",
			"accessed_count",
			"crypt_id",
		).
		Values(
			cred.ID,
			1,
			1,
			tagsJSON,
			0,
			cryptID,
		).
		RunWith(tx).
		ExecContext(ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	domainsJSON, err := mango.MarshalToString(cred.Domains)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	encryptedPwd, err := r.cipher.Encrypt(cred.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	detailsJSON, err := mango.MarshalToString(cred.Details)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	encryptedDetails, err := r.cipher.Encrypt(detailsJSON)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = sq.Insert("credential_versions").
		Columns(
			"credential_id",
			"version",
			"service",
			"domains",
			"email",
			"username",
			"s_password",
			"description",
			"s_details",
		).
		Values(
			cred.ID,
			1,
			cred.Service,
			domainsJSON,
			cred.Email,
			cred.Username,
			encryptedPwd,
			cred.Description,
			encryptedDetails,
		).
		RunWith(tx).
		ExecContext(ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = sq.Update("crypts").
		Set("updated_at", time.Now()).
		Set("total_active_credentials", sq.Expr("total_active_credentials + 1")).
		RunWith(tx).
		Where(sq.Eq{"id": cryptID}).
		ExecContext(ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	creds, err := r.QueryCredentials(ctx, repos.QueryCredentialsFilter{ID: cred.ID})
	if err != nil {
		return nil, err
	}
	if len(creds) != 1 {
		return nil, errors.New("failed to insert credential")
	}

	return creds[0], nil
}

/******************************************************************************/

func (r *DbRepo) getLatestCredentialVersion(ctx context.Context, id string) (int, error) {
	var latestVersion int
	err := sq.Select("latest_version").
		From("credentials").
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRowContext(ctx).
		Scan(&latestVersion)

	if err != nil {
		return 0, err
	}

	return latestVersion, nil
}

func (r *DbRepo) UpdateCredential(ctx context.Context, cryptID string, cred *repos.Credential) (*repos.Credential, error) {
	if cred.ID == "" {
		return nil, errors.New("credential id is required")
	}

	latestVersion, err := r.getLatestCredentialVersion(ctx, cred.ID)
	if err != nil {
		return nil, err
	}
	nextVersion := latestVersion + 1

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	tagsJSON, err := mango.MarshalToString(cred.Tags)
	if err != nil {
		return nil, err
	}

	_, err = sq.Update("credentials").
		Set("tags", tagsJSON).
		Set("current_version", nextVersion).
		Set("latest_version", nextVersion).
		RunWith(tx).
		Where(sq.Eq{"id": cred.ID}).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	domainsJSON, err := mango.MarshalToString(cred.Domains)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	encryptedPwd, err := r.cipher.Encrypt(cred.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	detailsJSON, err := mango.MarshalToString(cred.Details)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	encryptedDetails, err := r.cipher.Encrypt(detailsJSON)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = sq.Insert("credential_versions").
		Columns(
			"credential_id",
			"version",
			"service",
			"domains",
			"email",
			"username",
			"s_password",
			"description",
			"s_details",
		).
		Values(
			cred.ID,
			nextVersion,
			cred.Service,
			domainsJSON,
			cred.Email,
			cred.Username,
			encryptedPwd,
			cred.Description,
			encryptedDetails,
		).
		RunWith(tx).
		ExecContext(ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = sq.Update("crypts").
		Set("updated_at", time.Now()).
		RunWith(tx).
		Where(sq.Eq{"id": cryptID}).
		ExecContext(ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	creds, err := r.QueryCredentials(ctx, repos.QueryCredentialsFilter{ID: cred.ID})
	if err != nil {
		return nil, err
	}
	if len(creds) != 1 {
		return nil, errors.New("failed to update credential")
	}

	return creds[0], nil
}
