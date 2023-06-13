package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"account-service/internal/domain/secret"
	"account-service/pkg/store"
)

type SecretRepository struct {
	db map[string]secret.Entity
	sync.RWMutex
}

func NewSecretRepository() *SecretRepository {
	return &SecretRepository{
		db: make(map[string]secret.Entity),
	}
}

func (r *SecretRepository) Select(ctx context.Context) (dest []secret.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]secret.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *SecretRepository) Create(ctx context.Context, data secret.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *SecretRepository) Get(ctx context.Context, id string) (dest secret.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest, ok := r.db[id]
	if !ok {
		err = store.ErrorNotFound
		return
	}

	return
}

func (r *SecretRepository) Update(ctx context.Context, id string, data secret.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	r.db[id] = data

	return
}

func (r *SecretRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	delete(r.db, id)

	return
}

func (r *SecretRepository) generateID() string {
	return uuid.New().String()
}
