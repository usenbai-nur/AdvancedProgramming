package cars

import (
	"errors"
	"strings"
	"time"
)

var ErrValidation = errors.New("validation error")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req CreateCarRequest) (Car, error) {
	brand := strings.TrimSpace(req.Brand)
	model := strings.TrimSpace(req.Model)

	if brand == "" || model == "" {
		return Car{}, ErrValidation
	}
	if req.Year < 1950 || req.Year > time.Now().Year()+1 {
		return Car{}, ErrValidation
	}
	if req.Price <= 0 {
		return Car{}, ErrValidation
	}
	if req.Mileage < 0 {
		return Car{}, ErrValidation
	}

	car := Car{
		Brand:     brand,
		Model:     model,
		Year:      req.Year,
		Price:     req.Price,
		Mileage:   req.Mileage,
		Status:    StatusAvailable,
		CreatedAt: time.Now().UTC(),
	}

	created := s.repo.Create(car)
	return created, nil
}

func (s *Service) GetByID(id int) (Car, error) {
	return s.repo.GetByID(id)
}

func (s *Service) List() []Car {
	return s.repo.List()
}
