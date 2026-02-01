package orders

// Запрос на создание заказа (POST /api/orders)
type CreateOrderRequest struct {
	UserID  int    `json:"user_id" validate:"required,min=1"`
	CarID   int    `json:"car_id" validate:"required,min=1"`
	Comment string `json:"comment,omitempty"`
}

// Запрос на обновление статуса (PUT /api/orders/status?id=1)
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

// Унифицированный ответ API
type OrderResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
