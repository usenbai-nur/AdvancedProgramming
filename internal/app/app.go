package app

import (
	"fmt"
	"log"
	"net/http"

	"AdvancedProgramming/internal/auth"
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

	// Public endpoints
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			"Car Store API\nTeam: Nurdaulet, Nurbol, Ehson\n\nEndpoints:\n"+
				"- GET /health\n"+
				"- POST /auth/register\n"+
				"- POST /auth/login\n"+
				"- GET/POST /cars (protected)\n"+
				"- GET/PUT/DELETE /cars/{id} (protected)\n"+
				"- GET/POST /orders (protected)\n"+
				"- PUT /orders/status (protected)\n"+
				"\nUI:\n"+
				"- GET /ui/cars\n"+
				"- GET /ui/cars/new\n"+
				"- GET /ui/orders\n"+
				"- GET /ui/login\n"+
				"- GET /ui/register\n",
		)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Server is up and running!"))
	})

	// Auth (public)
	mux.HandleFunc("/auth/register", auth.Register)
	mux.HandleFunc("/auth/login", auth.Login)

	// Cars (protected API)
	carRepo := cars.NewRepository()
	carService := cars.NewService(carRepo)
	carHandler := cars.NewHandler(carService)

	// wrap handlers with middleware (http.Handler)
	mux.Handle("/cars", auth.AuthMiddleware(http.HandlerFunc(carHandler.Cars)))
	mux.Handle("/cars/", auth.AuthMiddleware(http.HandlerFunc(carHandler.CarByID)))

	// Web UI (public UI; if you want protect it too, we can wrap /ui/* later)
	webui.Register(mux, carService)

	// Orders (protected API)
	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(&orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	mux.Handle("/orders", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})))

	mux.Handle("/orders/status", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.UpdateOrderStatus(w, r)
	})))

	fmt.Println("Car Store API started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
