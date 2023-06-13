package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"account-service/internal/domain/user"
	"account-service/pkg/store"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (s *UsersRepository) Create(ctx context.Context, data user.Entity) (dest string, err error) {
	query := `
		INSERT INTO users (phone, type)
		VALUES ($1, $2)
		RETURNING id`

	args := []any{data.Phone, data.Type}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&dest)

	return
}

func (s *UsersRepository) GetByID(ctx context.Context, id string) (dest user.Entity, err error) {
	query := `
		SELECT created_at, updated_at, id, phone, type, name, email, birth_date
		FROM users
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

func (s *UsersRepository) GetByPhone(ctx context.Context, phone string) (dest user.Entity, err error) {
	query := `
		SELECT id
		FROM users
		WHERE phone=$1`

	args := []any{phone}

	if err = s.db.GetContext(ctx, &dest, query, args...); err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = store.ErrorNotFound
	}

	return
}

func (s *UsersRepository) Update(ctx context.Context, id string, data user.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")
		query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))

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

func (s *UsersRepository) prepareArgs(data user.Entity) (sets []string, args []any) {
	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	if data.Email != nil {
		args = append(args, data.Email)
		sets = append(sets, fmt.Sprintf("email=$%d", len(args)))
	}

	if data.BirthDate != nil {
		args = append(args, data.BirthDate)
		sets = append(sets, fmt.Sprintf("birth_date=$%d", len(args)))
	}

	return
}
