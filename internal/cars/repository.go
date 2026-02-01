package cars

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrNotFound = errors.New("car not found")

type Repository struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int]Car
}

func NewRepository() *Repository {
	return &Repository{
		items:  make(map[int]Car),
		nextID: 0,
	}
}

func (r *Repository) Create(c Car) Car {
	id := int(atomic.AddInt64(&r.nextID, 1))
	c.ID = id

	r.mu.Lock()
	r.items[id] = c
	r.mu.Unlock()

	return c
}

func (r *Repository) GetByID(id int) (Car, error) {
	r.mu.RLock()
	c, ok := r.items[id]
	r.mu.RUnlock()

	if !ok {
		return Car{}, ErrNotFound
	}
	return c, nil
}

func (r *Repository) List() []Car {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Car, 0, len(r.items))
	for _, c := range r.items {
		out = append(out, c)
	}
	return out
}
