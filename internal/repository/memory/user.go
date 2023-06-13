package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"account-service/internal/domain/user"
	"account-service/pkg/store"
)

type UsersRepository struct {
	db map[string]user.Entity
	sync.RWMutex
}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{
		db: make(map[string]user.Entity),
	}
}

func (r *UsersRepository) Select(ctx context.Context) (dest []user.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]user.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *UsersRepository) Create(ctx context.Context, data user.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *UsersRepository) GetByID(ctx context.Context, id string) (dest user.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[id]
	if !ok {
		err = store.ErrorNotFound
		return
	}

	return
}

func (r *UsersRepository) GetByPhone(ctx context.Context, phone string) (dest user.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[phone]
	if !ok {
		err = store.ErrorNotFound
		return
	}

	return
}

func (r *UsersRepository) Update(ctx context.Context, id string, data user.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	r.db[id] = data

	return
}

func (r *UsersRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	delete(r.db, id)

	return
}

func (r *UsersRepository) generateID() string {
	return uuid.New().String()
}
