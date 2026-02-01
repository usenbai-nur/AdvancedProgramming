package cars

type CreateCarRequest struct {
	Brand   string `json:"brand"`
	Model   string `json:"model"`
	Year    int    `json:"year"`
	Price   int    `json:"price"`
	Mileage int    `json:"mileage"`
}

type CarResponse struct {
	Car
}
