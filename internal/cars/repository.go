package cars

import (
	"errors"
	"sort"
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

	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (r *Repository) Update(id int, updateFn func(Car) (Car, error)) (Car, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.items[id]
	if !ok {
		return Car{}, ErrNotFound
	}

	updated, err := updateFn(current)
	if err != nil {
		return Car{}, err
	}

	// protect system fields
	updated.ID = id
	updated.CreatedAt = current.CreatedAt

	r.items[id] = updated
	return updated, nil
}

func (r *Repository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
