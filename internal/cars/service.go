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

func (s *Service) Update(id int, req UpdateCarRequest) (Car, error) {
	return s.repo.Update(id, func(current Car) (Car, error) {
		updated := current

		if req.Brand != nil {
			v := strings.TrimSpace(*req.Brand)
			if v == "" {
				return Car{}, ErrValidation
			}
			updated.Brand = v
		}

		if req.Model != nil {
			v := strings.TrimSpace(*req.Model)
			if v == "" {
				return Car{}, ErrValidation
			}
			updated.Model = v
		}

		if req.Year != nil {
			if *req.Year < 1950 || *req.Year > time.Now().Year()+1 {
				return Car{}, ErrValidation
			}
			updated.Year = *req.Year
		}

		if req.Price != nil {
			if *req.Price <= 0 {
				return Car{}, ErrValidation
			}
			updated.Price = *req.Price
		}

		if req.Mileage != nil {
			if *req.Mileage < 0 {
				return Car{}, ErrValidation
			}
			updated.Mileage = *req.Mileage
		}

		if req.Status != nil {
			switch *req.Status {
			case StatusAvailable, StatusReserved, StatusSold:
				updated.Status = *req.Status
			default:
				return Car{}, ErrValidation
			}
		}

		return updated, nil
	})
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
