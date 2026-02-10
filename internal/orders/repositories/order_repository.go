package repositories

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"time"

	"AdvancedProgramming/internal/infrastructure"
	"AdvancedProgramming/internal/orders/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	nextID int64
}

func NewOrderRepository() OrderRepository {
	return OrderRepository{nextID: 0}
}

func (r *OrderRepository) Create(order models.Order) (models.Order, error) {
	order.ID = int(atomic.AddInt64(&r.nextID, 1))
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if order.Status == "" {
		order.Status = "pending"
	}

	_, err := infrastructure.Database.Collection("orders").InsertOne(context.TODO(), order)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *OrderRepository) GetByID(id int) (models.Order, error) {
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
	var order models.Order
	err := infrastructure.Database.Collection("orders").FindOne(
		context.TODO(), bson.M{"id": id},
	).Decode(&order)
	if err != nil {
		return models.Order{}, errors.New("order not found")
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	_, err = infrastructure.Database.Collection("orders").ReplaceOne(
		context.TODO(), bson.M{"id": id}, order,
	)
	return order, err
}

func (r *OrderRepository) Delete(id int) error {
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

// НОВЫЕ МЕТОДЫ:

// GetByStatus - получить заказы по статусу
func (r *OrderRepository) GetByStatus(status string) ([]models.Order, error) {
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

// GetRecent - получить последние N заказов
func (r *OrderRepository) GetRecent(limit int) ([]models.Order, error) {
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

// Search - поиск по комментарию
func (r *OrderRepository) Search(query string) ([]models.Order, error) {
	var orders []models.Order

	// Получаем все заказы
	allOrders, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	// Фильтруем по комментарию (простой поиск)
	for _, order := range allOrders {
		if strings.Contains(strings.ToLower(order.Comment), query) {
			orders = append(orders, order)
		}
	}

	return orders, nil
}
