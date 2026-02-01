package repositories

import (
	"errors"
	"sync"
	"time"

	"AdvancedProgramming/models"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[int]*models.Order
	lastID int
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int]*models.Order),
	}
}

func (r *OrderRepository) Create(order *models.Order) (*models.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	order.ID = r.lastID
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	if order.Status == "" {
		order.Status = "pending"
	}

	r.orders[order.ID] = order
	return order, nil
}

func (r *OrderRepository) GetByID(id int) (*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (r *OrderRepository) GetAll() ([]*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*models.Order, 0, len(r.orders))
	for _, o := range r.orders {
		result = append(result, o)
	}
	return result, nil
}

func (r *OrderRepository) GetByUserID(userID int) ([]*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*models.Order
	for _, o := range r.orders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}
	return result, nil
}

func (r *OrderRepository) UpdateStatus(id int, status string) (*models.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	return order, nil
}
