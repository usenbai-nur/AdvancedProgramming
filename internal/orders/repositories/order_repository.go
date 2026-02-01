package repositories

import (
	"context"
	"errors"
	"time"

	"AdvancedProgramming/internal/infrastructure"
	"AdvancedProgramming/internal/orders/models"
	"go.mongodb.org/mongo-driver/bson"
)

type OrderRepository struct{}

func NewOrderRepository() OrderRepository { return OrderRepository{} }

func (r OrderRepository) Create(order models.Order) (models.Order, error) {
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

func (r OrderRepository) GetByID(id int) (models.Order, error) {
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

func (r OrderRepository) GetAll() ([]models.Order, error) {
	cursor, err := infrastructure.Database.Collection("orders").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	if err = cursor.All(context.TODO(), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r OrderRepository) GetByUserID(userID int) ([]models.Order, error) {
	cursor, err := infrastructure.Database.Collection("orders").Find(
		context.TODO(),
		bson.M{"userid": userID},
	)
	if err != nil {
		return nil, err
	}
	var orders []models.Order
	return orders, cursor.All(context.TODO(), &orders)
}

func (r OrderRepository) UpdateStatus(id int, status string) (models.Order, error) {
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
