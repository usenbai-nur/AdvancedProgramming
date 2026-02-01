package handlers

type CreateOrderRequest struct {
	UserID  int    `json:"user_id"`
	CarID   int    `json:"car_id"`
	Comment string `json:"comment"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
