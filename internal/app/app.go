package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		log.Printf("Database init warning: %v", err)
	}
	defer infrastructure.CloseDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			"Car Store API\nTeam: Nurdaulet, Nurbol, Ehson\n\n"+
				"=== Endpoints ===\n\n"+
				"Health:\n"+
				"  GET  /health\n\n"+
				"Auth:\n"+
				"  POST /auth/register\n"+
				"  POST /auth/login\n"+
				"  GET  /auth/me\n"+
				"  POST /auth/favorites/{carID}\n\n"+
				"Cars:\n"+
				"  GET    /cars\n"+
				"  GET    /cars/{id}\n"+
				"  POST   /cars              (admin)\n"+
				"  PUT    /cars/{id}         (admin)\n"+
				"  DELETE /cars/{id}         (admin)\n\n"+
				"Orders:\n"+
				"  POST   /orders            (user/admin)\n"+
				"  GET    /orders            (admin)\n"+
				"  GET    /orders/{id}       (admin)\n"+
				"  PUT    /orders/{id}       (admin)\n"+
				"  DELETE /orders/{id}       (admin)\n"+
				"  GET    /users/{id}/orders (admin)\n"+
				"  GET    /orders/stats      (admin)\n"+
				"  GET    /orders/search?q=  (admin)\n\n"+
				"UI:\n"+
				"  GET /ui/cars\n"+
				"  GET /ui/cars/new\n"+
				"  GET /ui/orders\n"+
				"  GET /ui/login\n"+
				"  GET /ui/register\n",
		)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("✅ Server is up and running!\n"))
	})

	mux.HandleFunc("/auth/register", auth.Register)
	mux.HandleFunc("/auth/login", auth.Login)
	mux.Handle("/auth/me", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		username, ok := auth.UsernameFromContext(r.Context())
		if !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		user, ok := auth.GetUserByUsername(username)
		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"user": user})
	})))

	mux.Handle("/auth/favorites/", auth.RequireRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		username, _ := auth.UsernameFromContext(r.Context())
		idStr := strings.TrimPrefix(r.URL.Path, "/auth/favorites/")
		carID, err := strconv.Atoi(idStr)
		if err != nil || carID <= 0 {
			http.Error(w, "invalid car id", http.StatusBadRequest)
			return
		}
		user, err := auth.AddFavorite(username, carID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "favorite added", "user": user})
	}), auth.RoleUser, auth.RoleAdmin))

	carRepo := cars.NewRepository()
	carService := cars.NewService(carRepo)
	carHandler := cars.NewHandler(carService)
	mux.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			auth.RequireRoles(http.HandlerFunc(carHandler.Cars), auth.RoleAdmin).ServeHTTP(w, r)
			return
		}
		carHandler.Cars(w, r)
	})
	mux.HandleFunc("/cars/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut || r.Method == http.MethodDelete {
			auth.RequireRoles(http.HandlerFunc(carHandler.CarByID), auth.RoleAdmin).ServeHTTP(w, r)
			return
		}
		carHandler.CarByID(w, r)
	})

	webui.Register(mux, carService)

	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService(&orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	mux.Handle("/orders/stats", auth.RequireRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.GetOrderStats(w, r)
	}), auth.RoleAdmin))

	mux.Handle("/orders/search", auth.RequireRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.SearchOrders(w, r)
	}), auth.RoleAdmin))

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			auth.RequireRoles(http.HandlerFunc(orderHandler.CreateOrder), auth.RoleUser, auth.RoleAdmin).ServeHTTP(w, r)
		case http.MethodGet:
			auth.RequireRoles(http.HandlerFunc(orderHandler.GetAllOrders), auth.RoleAdmin).ServeHTTP(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/orders/", auth.RequireRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPut && r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.HandleOrderByID(w, r)
	}), auth.RoleAdmin))

	mux.Handle("/users/", auth.RequireRoles(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		orderHandler.GetUserOrders(w, r)
	}), auth.RoleAdmin))

	fmt.Println("Car Store API started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
