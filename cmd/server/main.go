package main

import (
	"AdvancedProgramming/internal/infrastructure"
	"AdvancedProgramming/internal/orders/handlers"
	"AdvancedProgramming/internal/orders/repositories"
	"AdvancedProgramming/internal/orders/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	if err := infrastructure.InitDatabase(); err != nil {
		log.Fatal(" MongoDB Atlas failed:", err)
	}
	defer infrastructure.CloseDatabase()

	mux := http.NewServeMux()

	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(&orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, " Car Store API + MongoDB Atlas \nTeam: Nurdaulet, Nurbol, Ehson\nhttp://localhost:8080/orders")
	})

	mux.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Cars Module - Nurdaulet (soon!)")
	})

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Auth Module - Ehson (soon!)")
	})

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodGet:
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

	mux.HandleFunc("/orders/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.UpdateOrderStatus(w, r)
	})

	fmt.Println("Car Store API + MongoDB Atlas started!")
	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
