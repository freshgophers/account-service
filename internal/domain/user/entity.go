package user

import (
	"time"
)

const (
	CUSTOMER = "customer"
)

type Entity struct {
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	ID        string     `db:"id"`
	Phone     string     `db:"phone"`
	Type      string     `db:"type"`
	Name      *string    `db:"name"`
	Email     *string    `db:"email"`
	BirthDate *time.Time `db:"birth_date"`
}
