package services

import (
	"AdvancedProgramming/internal/orders/models"
	"AdvancedProgramming/internal/orders/repositories"
	"errors"
	"log"
	"strings"
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
	log.Println("📦 Order background processor started")
	for orderID := range s.processChan {
		log.Printf("⏳ Processing order %d ...", orderID)
		time.Sleep(3 * time.Second)

		_, err := s.repo.UpdateStatus(orderID, "confirmed")
		if err != nil {
			log.Printf("❌ Failed to auto-confirm order %d: %v", orderID, err)
		} else {
			log.Printf("✅ Order %d automatically confirmed", orderID)
		}
	}
}

func (s *OrderService) CreateOrder(userID, carID int, comment string) (models.Order, error) {
	if userID <= 0 {
		return models.Order{}, errors.New("user_id must be positive")
	}
	if carID <= 0 {
		return models.Order{}, errors.New("car_id must be positive")
	}

	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return models.Order{}, errors.New("comment cannot be empty")
	}
	if len(comment) > 500 {
		return models.Order{}, errors.New("comment too long (max 500 characters)")
	}

	order := models.Order{
		UserID:  userID,
		CarID:   carID,
		Comment: comment,
		Status:  "pending",
	}

	created, err := s.repo.Create(order)
	if err != nil {
		return models.Order{}, err
	}

	go func() { s.processChan <- created.ID }()
	return created, nil
}

func (s *OrderService) GetOrder(id int) (models.Order, error) {
	if id <= 0 {
		return models.Order{}, errors.New("invalid order id")
	}
	return s.repo.GetByID(id)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.repo.GetAll()
}

func (s *OrderService) GetUserOrders(userID int) ([]models.Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetByUserID(userID)
}

func (s *OrderService) UpdateStatus(id int, status string) (models.Order, error) {
	if id <= 0 {
		return models.Order{}, errors.New("invalid order id")
	}
	if !s.validStatuses[status] {
		return models.Order{}, errors.New("invalid status. allowed: pending, confirmed, cancelled, completed")
	}
	return s.repo.UpdateStatus(id, status)
}

func (s *OrderService) DeleteOrder(id int) error {
	if id <= 0 {
		return errors.New("invalid order id")
	}
	return s.repo.Delete(id)
}

// НОВЫЕ МЕТОДЫ ДЛЯ УЛУЧШЕНИЯ:

// GetOrdersByStatus - фильтр по статусу
func (s *OrderService) GetOrdersByStatus(status string) ([]models.Order, error) {
	if !s.validStatuses[status] {
		return nil, errors.New("invalid status")
	}
	return s.repo.GetByStatus(status)
}

// GetRecentOrders - последние N заказов
func (s *OrderService) GetRecentOrders(limit int) ([]models.Order, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return s.repo.GetRecent(limit)
}

// GetOrderStats - статистика по заказам
func (s *OrderService) GetOrderStats() (map[string]interface{}, error) {
	allOrders, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total":     len(allOrders),
		"pending":   0,
		"confirmed": 0,
		"cancelled": 0,
		"completed": 0,
	}

	// Подсчёт по статусам
	for _, order := range allOrders {
		if count, ok := stats[order.Status].(int); ok {
			stats[order.Status] = count + 1
		}
	}

	// Популярные машины (top 5)
	carCount := make(map[int]int)
	for _, order := range allOrders {
		carCount[order.CarID]++
	}

	type carStat struct {
		CarID int
		Count int
	}
	var topCars []carStat
	for carID, count := range carCount {
		topCars = append(topCars, carStat{CarID: carID, Count: count})
	}

	// Сортировка по количеству (простая)
	for i := 0; i < len(topCars)-1; i++ {
		for j := i + 1; j < len(topCars); j++ {
			if topCars[j].Count > topCars[i].Count {
				topCars[i], topCars[j] = topCars[j], topCars[i]
			}
		}
	}

	if len(topCars) > 5 {
		topCars = topCars[:5]
	}
	stats["top_cars"] = topCars

	// Заказы за сегодня
	today := time.Now().Truncate(24 * time.Hour)
	todayCount := 0
	for _, order := range allOrders {
		if order.CreatedAt.After(today) {
			todayCount++
		}
	}
	stats["today"] = todayCount

	return stats, nil
}

// SearchOrders - поиск по комментарию
func (s *OrderService) SearchOrders(query string) ([]models.Order, error) {
	if len(query) < 2 {
		return nil, errors.New("search query too short (min 2 characters)")
	}
	return s.repo.Search(strings.ToLower(query))
}
