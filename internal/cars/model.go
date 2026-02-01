package cars

import "time"

type Status string

const (
	StatusAvailable Status = "available"
	StatusReserved  Status = "reserved"
	StatusSold      Status = "sold"
)

type Car struct {
	ID        int       `json:"id"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Year      int       `json:"year"`
	Price     int       `json:"price"`
	Mileage   int       `json:"mileage"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
