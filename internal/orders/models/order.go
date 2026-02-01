package models

import "time"

type Order struct {
	ID        int       `json:"id" bson:"id"`
	UserID    int       `json:"userid" bson:"userid"`
	CarID     int       `json:"carid" bson:"carid"`
	Status    string    `json:"status" bson:"status"`
	Comment   string    `json:"comment" bson:"comment"`
	CreatedAt time.Time `json:"createdat" bson:"createdat"`
	UpdatedAt time.Time `json:"updatedat" bson:"updatedat"`
}

type OrderWithDetails struct {
	Order
	CarBrand string  `json:"car_brand"`
	CarModel string  `json:"car_model"`
	CarPrice float64 `json:"car_price"`
}
