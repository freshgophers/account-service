package secret

import "context"

type Repository interface {
	Create(ctx context.Context, data Entity) (id string, err error)
	GetByID(ctx context.Context, id string) (dest Entity, err error)
	Update(ctx context.Context, id string, data Entity) (err error)
}
