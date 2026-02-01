package cars

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"AdvancedProgramming/internal/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /cars  |  GET /cars
func (h *Handler) Cars(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := h.svc.List()
		httpx.WriteJSON(w, http.StatusOK, list)
		return

	case http.MethodPost:
		var req CreateCarRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, httpx.Err("bad_json", "Invalid JSON body"))
			return
		}

		created, err := h.svc.Create(req)
		if err != nil {
			if err == ErrValidation {
				httpx.WriteError(w, http.StatusBadRequest, httpx.Err("validation_error", "Invalid car fields"))
				return
			}
			httpx.WriteError(w, http.StatusInternalServerError, httpx.Err("server_error", "Internal server error"))
			return
		}

		httpx.WriteJSON(w, http.StatusCreated, created)
		return

	default:
		httpx.WriteError(w, http.StatusMethodNotAllowed, httpx.Err("method_not_allowed", "Method not allowed"))
		return
	}
}

// GET /cars/{id}
func (h *Handler) CarByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteError(w, http.StatusMethodNotAllowed, httpx.Err("method_not_allowed", "Method not allowed"))
		return
	}

	// r.URL.Path example: /cars/123
	path := strings.TrimPrefix(r.URL.Path, "/cars/")
	if path == "" || strings.Contains(path, "/") {
		httpx.WriteError(w, http.StatusNotFound, httpx.Err("not_found", "Not found"))
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil || id <= 0 {
		httpx.WriteError(w, http.StatusBadRequest, httpx.Err("bad_id", "Invalid car id"))
		return
	}

	car, err := h.svc.GetByID(id)
	if err != nil {
		if err == ErrNotFound {
			httpx.WriteError(w, http.StatusNotFound, httpx.Err("not_found", "Car not found"))
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, httpx.Err("server_error", "Internal server error"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, car)
}
