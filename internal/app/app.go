package app

import (
	"car-store/internal/infrastructure"
	"fmt"
	"net/http"
	// "car-store/internal/cars"
	// "car-store/internal/orders"
)

func Run() {

	config := infrastructure.LoadConfig()

	// carsRepo := cars.NewRepository()
	// ordersRepo := orders.NewRepository()

	// orderChan := make(chan string) // Можно заменить на тип Order позже

	// go orders.StartWorker(orderChan)

	mux := http.NewServeMux()

	// cars.RegisterRoutes(mux, carsRepo)
	// orders.RegisterRoutes(mux, ordersRepo, orderChan)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is up and running!"))
	})

	fmt.Printf("Сервер запущен на порту %s\n", config.Port)

	if err := http.ListenAndServe(config.Port, mux); err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err)
	}
}
