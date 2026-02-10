package app

import (
	"AdvancedProgramming/internal/auth"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", auth.Register)
	mux.HandleFunc("/login", auth.Login)

	pingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("username")
		fmt.Fprintf(w, "Status: authorized, User: %v", user)
	})

	mux.Handle("/api/ping", auth.AuthMiddleware(pingHandler))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
