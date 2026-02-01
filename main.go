package main

import (
	"fmt"
	"log"
	"net/http"

	"FinalProject/internal/orders/handlers"
	"FinalProject/internal/orders/repositories"
	"FinalProject/internal/orders/services"
)

func main() {
	mux := http.NewServeMux()

	// Orders module
	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Main page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Car Store - Team: Nurdaulet, Nurbol, Ehson")
	})

	mux.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Cars Module - Nurdaulet")
	})

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Auth Module - Ehson")
	})

	// ===== ORDERS API =====

	// /orders:
	// POST  /orders                – create order (JSON)
	// GET   /orders?id=1           – get one order
	// GET   /orders?user_id=1      – user's orders
	// GET   /orders                – all orders
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodGet:
			// priority: user_id → id → все
			if r.URL.Query().Get("user_id") != "" {
				orderHandler.GetUserOrders(w, r)
				return
			}
			if r.URL.Query().Get("id") != "" {
				orderHandler.GetOrder(w, r)
				return
			}
			orderHandler.GetAllOrders(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// PUT /orders/status?id=1 – update status
	mux.HandleFunc("/orders/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.UpdateOrderStatus(w, r)
	})

	// ===== Start server =====

	fmt.Println("Car Store API - http://localhost:8080")
	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
