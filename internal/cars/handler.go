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
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
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

// GET/PUT/DELETE /cars/{id}
func (h *Handler) CarByID(w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodGet:
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
		return

	case http.MethodPut:
		var req UpdateCarRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, httpx.Err("bad_json", "Invalid JSON body"))
			return
		}

		updated, err := h.svc.Update(id, req)
		if err != nil {
			if err == ErrNotFound {
				httpx.WriteError(w, http.StatusNotFound, httpx.Err("not_found", "Car not found"))
				return
			}
			if err == ErrValidation {
				httpx.WriteError(w, http.StatusBadRequest, httpx.Err("validation_error", "Invalid car fields"))
				return
			}
			httpx.WriteError(w, http.StatusInternalServerError, httpx.Err("server_error", "Internal server error"))
			return
		}

		httpx.WriteJSON(w, http.StatusOK, updated)
		return

	case http.MethodDelete:
		if err := h.svc.Delete(id); err != nil {
			if err == ErrNotFound {
				httpx.WriteError(w, http.StatusNotFound, httpx.Err("not_found", "Car not found"))
				return
			}
			httpx.WriteError(w, http.StatusInternalServerError, httpx.Err("server_error", "Internal server error"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return

	default:
		httpx.WriteError(w, http.StatusMethodNotAllowed, httpx.Err("method_not_allowed", "Method not allowed"))
		return
	}
}
