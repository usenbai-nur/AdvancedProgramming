package app

import (
	"fmt"
	"log"
	"net/http"

	"AdvancedProgramming/internal/cars"
	"AdvancedProgramming/internal/infrastructure"

	// Orders (Nurbol)
	"AdvancedProgramming/internal/orders/handlers"
	"AdvancedProgramming/internal/orders/repositories"
	"AdvancedProgramming/internal/orders/services"
)

func Run() {
	if err := infrastructure.InitDatabase(); err != nil {
		log.Fatal("Database init failed:", err)
	}
	defer infrastructure.CloseDatabase()

	mux := http.NewServeMux()

	// Health + Home
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			"Car Store API\nTeam: Nurdaulet, Nurbol, Ehson\n\n"+
				"=== Endpoints ===\n\n"+
				"Health:\n"+
				"  GET  /health\n\n"+
				"Cars (Nurdaulet):\n"+
				"  POST   /cars\n"+
				"  GET    /cars\n"+
				"  GET    /cars/{id}\n"+
				"  PUT    /cars/{id}\n"+
				"  DELETE /cars/{id}\n\n"+
				"Orders (Nurbol):\n"+
				"  POST   /orders\n"+
				"  GET    /orders\n"+
				"  GET    /orders/{id}\n"+
				"  PUT    /orders/{id}\n"+
				"  DELETE /orders/{id}\n"+
				"  GET    /users/{id}/orders\n"+
				"  GET    /orders/stats\n"+
				"  GET    /orders/search?q=query\n",
		)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("✅ Server is up and running!"))
	})

	// Cars (Nurdaulet) - in-memory storage
	carRepo := cars.NewRepository()
	carService := cars.NewService(carRepo)
	carHandler := cars.NewHandler(carService)
	cars.RegisterRoutes(mux, carHandler)

	// Orders (Nurbol) - MongoDB + RESTful endpoints
	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(&orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// ВАЖНО: Регистрируем специфичные роуты ПЕРВЫМИ!

	// GET /orders/stats - статистика
	mux.HandleFunc("/orders/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.GetOrderStats(w, r)
	})

	// GET /orders/search?q=query - поиск
	mux.HandleFunc("/orders/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.SearchOrders(w, r)
	})

	// POST /orders, GET /orders - базовая коллекция
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodGet:
			orderHandler.GetAllOrders(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET/PUT/DELETE /orders/{id} - конкретный заказ (ПОСЛЕ специфичных!)
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		orderHandler.HandleOrderByID(w, r)
	})

	// GET /users/{userId}/orders - заказы пользователя
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.GetUserOrders(w, r)
	})

	fmt.Println(" Car Store API started at http://localhost:8080")
	fmt.Println(" Orders Stats: http://localhost:8080/orders/stats")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
