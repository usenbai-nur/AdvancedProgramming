package services

import (
	"FinalProject/internal/orders/models"
	"FinalProject/internal/orders/repositories"
	"errors"
	"log"
	"time"
)

type OrderService struct {
	repo          *repositories.OrderRepository
	processChan   chan int
	validStatuses map[string]bool
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	s := &OrderService{
		repo:        repo,
		processChan: make(chan int, 10),
		validStatuses: map[string]bool{
			"pending":   true,
			"confirmed": true,
			"cancelled": true,
			"completed": true,
		},
	}

	go s.backgroundProcessor()

	return s
}

func (s *OrderService) backgroundProcessor() {
	log.Println("Order background processor started")
	for orderID := range s.processChan {
		log.Printf("Processing order %d ...", orderID)
		time.Sleep(2 * time.Second)
		log.Printf("Order %d processed", orderID)
	}
}

func (s *OrderService) CreateOrder(userID, carID int, comment string) (*models.Order, error) {
	if userID <= 0 || carID <= 0 {
		return nil, errors.New("invalid user_id or car_id")
	}

	order := &models.Order{
		UserID:  userID,
		CarID:   carID,
		Comment: comment,
		Status:  "pending",
	}

	created, err := s.repo.Create(order)
	if err != nil {
		return nil, err
	}

	go func() {
		s.processChan <- created.ID
	}()

	return created, nil
}

func (s *OrderService) GetOrder(id int) (*models.Order, error) {
	return s.repo.GetByID(id)
}

func (s *OrderService) GetAllOrders() ([]*models.Order, error) {
	return s.repo.GetAll()
}

func (s *OrderService) GetUserOrders(userID int) ([]*models.Order, error) {
	return s.repo.GetByUserID(userID)
}

func (s *OrderService) UpdateStatus(id int, status string) (*models.Order, error) {
	if !s.validStatuses[status] {
		return nil, errors.New("invalid status")
	}
	return s.repo.UpdateStatus(id, status)
}
