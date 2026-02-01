package cars

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	// collection routes
	mux.HandleFunc("/cars", h.Cars)

	// item routes
	mux.HandleFunc("/cars/", h.CarByID)
}
