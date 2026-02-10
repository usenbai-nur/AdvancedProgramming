package cars

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/cars", h.Cars)
	mux.HandleFunc("/cars/", h.CarByID)
}
