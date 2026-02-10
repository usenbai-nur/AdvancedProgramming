package app

import (
	"fmt"
	"log"
	"net/http"

	"AdvancedProgramming/internal/cars"
	"AdvancedProgramming/internal/infrastructure"
	"AdvancedProgramming/internal/orders/handlers"
	"AdvancedProgramming/internal/orders/repositories"
	"AdvancedProgramming/internal/orders/services"
	"AdvancedProgramming/internal/webui"
)

func Run() {
	if err := infrastructure.InitDatabase(); err != nil {
		log.Fatal("Database init failed:", err)
	}
	defer infrastructure.CloseDatabase()

	mux := http.NewServeMux()

	// Home
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			"Car Store API\nTeam: Nurdaulet, Nurbol, Ehson\n\nEndpoints:\n"+
				"- GET /health\n"+
				"- GET/POST /cars\n"+
				"- GET/PUT/DELETE /cars/{id}\n"+
				"- GET/POST /orders\n"+
				"- PUT /orders/status\n"+
				"\nUI:\n"+
				"- GET /ui/cars\n"+
				"- GET /ui/cars/new\n",
		)
	})

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Server is up and running!"))
	})

	// Cars (Nurdaulet) - in-memory storage
	carRepo := cars.NewRepository()
	carService := cars.NewService(carRepo)
	carHandler := cars.NewHandler(carService)
	cars.RegisterRoutes(mux, carHandler)

	// Web UI (templates + css)
	webui.Register(mux, carService)

	// Orders (Nurbol)
	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(&orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

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

	fmt.Println("Car Store API started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
