package models

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CarID     int       `json:"car_id"`
	Status    string    `json:"status"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderWithDetails struct {
	Order
	CarBrand string  `json:"car_brand"`
	CarModel string  `json:"car_model"`
	CarPrice float64 `json:"car_price"`
}
