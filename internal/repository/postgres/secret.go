package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"account-service/internal/domain/secret"
	"account-service/pkg/store"
)

type SecretRepository struct {
	db *sqlx.DB
}

func NewSecretRepository(db *sqlx.DB) *SecretRepository {
	return &SecretRepository{
		db: db,
	}
}

func (s *SecretRepository) Create(ctx context.Context, data secret.Entity) (id string, err error) {
	query := `
		INSERT INTO secrets (secret, phone, status)
		VALUES ($1, $2, $3) 
		RETURNING id`

	args := []any{data.Secret, data.Phone, data.Status}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *SecretRepository) Get(ctx context.Context, id string) (dest secret.Entity, err error) {
	query := `
		SELECT created_at, updated_at, id, secret, phone, attempts, status
		FROM secrets
		WHERE id=$1`

	args := []any{id}

	if err = s.db.GetContext(ctx, &dest, query, args...); err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = store.ErrorNotFound
	}

	return
}

func (s *SecretRepository) Update(ctx context.Context, id string, data secret.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")
		query := fmt.Sprintf("UPDATE secrets SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		fmt.Println(query)

		_, err = s.db.ExecContext(ctx, query, args...)
		if err != nil && err != sql.ErrNoRows {
			return
		}

		if err == sql.ErrNoRows {
			err = store.ErrorNotFound
		}
	}

	return
}

func (s *SecretRepository) prepareArgs(data secret.Entity) (sets []string, args []any) {
	if data.Attempts != nil {
		args = append(args, data.Attempts)
		sets = append(sets, fmt.Sprintf("attempts=$%d", len(args)))
	}

	if data.Status != nil {
		args = append(args, data.Status)
		sets = append(sets, fmt.Sprintf("status=$%d", len(args)))
	}

	return
}
