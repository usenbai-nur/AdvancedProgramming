package main

import (
	"fmt"
	"log"
	"net/http"

	"AdvancedProgramming/handlers"
	"AdvancedProgramming/repositories"
	"AdvancedProgramming/services"
)

func main() {
	mux := http.NewServeMux()

	// ==== Инициализация модулей ====

	// Orders module (твоя часть)
	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// ==== Роуты ====

	// Главная страница (как было)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Car Store - Team: Nurdaulet, Nurbol, Ehson")
	})

	// Cars – пока просто заглушка (Нурдаулет)
	mux.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Cars Module - Nurdaulet")
	})

	// Auth – заглушка (Эхсон)
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Auth Module - Ehson")
	})

	// ===== ORDERS API (твоя реальная логика) =====

	// /orders:
	// POST  /orders                – создать заказ (JSON)
	// GET   /orders?id=1           – получить один заказ
	// GET   /orders?user_id=1      – заказы пользователя
	// GET   /orders                – все заказы
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodGet:
			// приоритет: user_id → id → все
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

	// PUT /orders/status?id=1 – обновить статус
	mux.HandleFunc("/orders/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.UpdateOrderStatus(w, r)
	})

	// ===== Запуск сервера =====

	fmt.Println("Car Store API - http://localhost:8080")
	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
