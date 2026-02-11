package repositories

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"AdvancedProgramming/internal/infrastructure"
	"AdvancedProgramming/internal/orders/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	nextID int64
	mu     sync.RWMutex
	items  map[int]models.Order
}

func NewOrderRepository() OrderRepository {
	return OrderRepository{nextID: 0, items: make(map[int]models.Order)}
}

func (r *OrderRepository) useMemory() bool {
	return infrastructure.Database == nil
}

func (r *OrderRepository) Create(order models.Order) (models.Order, error) {
	order.ID = int(atomic.AddInt64(&r.nextID, 1))
	order.CreatedAt = time.Now().UTC()
	order.UpdatedAt = order.CreatedAt
	if order.Status == "" {
		order.Status = "pending"
	}

	if r.useMemory() {
		r.mu.Lock()
		r.items[order.ID] = order
		r.mu.Unlock()
		return order, nil
	}

	_, err := infrastructure.Database.Collection("orders").InsertOne(context.TODO(), order)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *OrderRepository) GetByID(id int) (models.Order, error) {
	if r.useMemory() {
		r.mu.RLock()
		order, ok := r.items[id]
		r.mu.RUnlock()
		if !ok {
			return models.Order{}, errors.New("order not found")
		}
		return order, nil
	}

	var order models.Order
	err := infrastructure.Database.Collection("orders").FindOne(
		context.TODO(),
		bson.M{"id": id},
	).Decode(&order)
	if err != nil {
		return models.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *OrderRepository) GetAll() ([]models.Order, error) {
	if r.useMemory() {
		r.mu.RLock()
		orders := make([]models.Order, 0, len(r.items))
		for _, order := range r.items {
			orders = append(orders, order)
		}
		r.mu.RUnlock()
		sortOrdersByCreatedDesc(orders)
		return orders, nil
	}

	cursor, err := infrastructure.Database.Collection("orders").Find(
		context.TODO(),
		bson.M{},
		options.Find().SetSort(bson.M{"createdat": -1}),
	)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	if err = cursor.All(context.TODO(), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetByUserID(userID int) ([]models.Order, error) {
	if r.useMemory() {
		r.mu.RLock()
		orders := make([]models.Order, 0)
		for _, order := range r.items {
			if order.UserID == userID {
				orders = append(orders, order)
			}
		}
		r.mu.RUnlock()
		sortOrdersByCreatedDesc(orders)
		return orders, nil
	}

	cursor, err := infrastructure.Database.Collection("orders").Find(
		context.TODO(),
		bson.M{"userid": userID},
		options.Find().SetSort(bson.M{"createdat": -1}),
	)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	return orders, cursor.All(context.TODO(), &orders)
}

func (r *OrderRepository) UpdateStatus(id int, status string) (models.Order, error) {
	if r.useMemory() {
		r.mu.Lock()
		order, ok := r.items[id]
		if !ok {
			r.mu.Unlock()
			return models.Order{}, errors.New("order not found")
		}
		order.Status = status
		order.UpdatedAt = time.Now().UTC()
		r.items[id] = order
		r.mu.Unlock()
		return order, nil
	}

	var order models.Order
	err := infrastructure.Database.Collection("orders").FindOne(
		context.TODO(), bson.M{"id": id},
	).Decode(&order)
	if err != nil {
		return models.Order{}, errors.New("order not found")
	}

	order.Status = status
	order.UpdatedAt = time.Now().UTC()
	_, err = infrastructure.Database.Collection("orders").ReplaceOne(
		context.TODO(), bson.M{"id": id}, order,
	)
	return order, err
}

func (r *OrderRepository) Delete(id int) error {
	if r.useMemory() {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, ok := r.items[id]; !ok {
			return errors.New("order not found")
		}
		delete(r.items, id)
		return nil
	}

	result, err := infrastructure.Database.Collection("orders").DeleteOne(
		context.TODO(),
		bson.M{"id": id},
	)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (r *OrderRepository) GetByStatus(status string) ([]models.Order, error) {
	if r.useMemory() {
		r.mu.RLock()
		orders := make([]models.Order, 0)
		for _, order := range r.items {
			if order.Status == status {
				orders = append(orders, order)
			}
		}
		r.mu.RUnlock()
		sortOrdersByCreatedDesc(orders)
		return orders, nil
	}

	cursor, err := infrastructure.Database.Collection("orders").Find(
		context.TODO(),
		bson.M{"status": status},
		options.Find().SetSort(bson.M{"createdat": -1}),
	)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	return orders, cursor.All(context.TODO(), &orders)
}

func (r *OrderRepository) GetRecent(limit int) ([]models.Order, error) {
	if r.useMemory() {
		all, err := r.GetAll()
		if err != nil {
			return nil, err
		}
		if len(all) > limit {
			all = all[:limit]
		}
		return all, nil
	}

	cursor, err := infrastructure.Database.Collection("orders").Find(
		context.TODO(),
		bson.M{},
		options.Find().SetSort(bson.M{"createdat": -1}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	return orders, cursor.All(context.TODO(), &orders)
}

func (r *OrderRepository) Search(query string) ([]models.Order, error) {
	allOrders, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	for _, order := range allOrders {
		if strings.Contains(strings.ToLower(order.Comment), query) {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func sortOrdersByCreatedDesc(orders []models.Order) {
	for i := 0; i < len(orders)-1; i++ {
		for j := i + 1; j < len(orders); j++ {
			if orders[j].CreatedAt.After(orders[i].CreatedAt) {
				orders[i], orders[j] = orders[j], orders[i]
			}
		}
	}
}
