package handlers

import (
	"AdvancedProgramming/internal/orders/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(s *services.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

// CreateOrder - POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "invalid JSON body",
		})
		return
	}

	order, err := h.service.CreateOrder(req.UserID, req.CarID, req.Comment)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: "order created successfully",
		Data:    order,
	})
}

// GetAllOrders - GET /orders
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    orders,
	})
}

// HandleOrderByID - GET/PUT/DELETE /orders/{id}
func (h *OrderHandler) HandleOrderByID(w http.ResponseWriter, r *http.Request) {
	// Парсим ID из пути /orders/{id}
	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	if path == "" || strings.Contains(path, "/") {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: "not found",
		})
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "invalid order id",
		})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getOrderByID(w, id)
	case http.MethodPut:
		h.updateOrderStatus(w, r, id)
	case http.MethodDelete:
		h.deleteOrder(w, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// getOrderByID - GET /orders/{id}
func (h *OrderHandler) getOrderByID(w http.ResponseWriter, id int) {
	order, err := h.service.GetOrder(id)
	if err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    order,
	})
}

// updateOrderStatus - PUT /orders/{id}
func (h *OrderHandler) updateOrderStatus(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "invalid JSON body",
		})
		return
	}

	order, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "status updated successfully",
		Data:    order,
	})
}

// deleteOrder - DELETE /orders/{id}
func (h *OrderHandler) deleteOrder(w http.ResponseWriter, id int) {
	err := h.service.DeleteOrder(id)
	if err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "order deleted successfully",
	})
}

// GetUserOrders - GET /users/{userId}/orders
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	// Парсим /users/{userId}/orders
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[1] != "orders" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "invalid path, expected /users/{id}/orders",
		})
		return
	}

	userID, err := strconv.Atoi(parts[0])
	if err != nil || userID <= 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "invalid user id",
		})
		return
	}

	orders, err := h.service.GetUserOrders(userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    orders,
	})
}

// GetOrderStats - GET /orders/stats
func (h *OrderHandler) GetOrderStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetOrderStats()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

// SearchOrders - GET /orders/search?q=query
func (h *OrderHandler) SearchOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "search query required (?q=...)",
		})
		return
	}

	orders, err := h.service.SearchOrders(query)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    orders,
	})
}

func respondJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
