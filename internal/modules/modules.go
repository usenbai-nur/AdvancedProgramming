package models

import "time"

// User представляет пользователя системы (Админ или Клиент)
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`    // Пароль скрыт при отправке JSON
	Role     string `json:"role"` // Admin или User
}

// Car представляет автомобиль, доступный для аренды
type Car struct {
	ID          int     `json:"id"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	PricePerDay float64 `json:"price_per_day"`
	Status      string  `json:"status"` // например: "available", "booked"
}

// Order представляет транзакцию бронирования
type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CarID      int       `json:"car_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
}
