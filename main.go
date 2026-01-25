package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/cars", carsPage)
	http.HandleFunc("/orders", ordersPage)
	http.HandleFunc("/auth", authPage)

	fmt.Println("Car Store API - http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Car Store - Team: Nurdaulet, Nurbol, Ehson")
}

func carsPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Cars Module - Nurdaulet")
}

func ordersPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Orders Module - Nurbol")
}

func authPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth Module - Ehson")
}
