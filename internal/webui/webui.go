package webui

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"AdvancedProgramming/internal/cars"
)

type Handler struct {
	cars *cars.Service
	tmpl *template.Template
}

func Register(mux *http.ServeMux, carService *cars.Service) {
	// Static files (CSS)
	mux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))),
	)

	h := &Handler{
		cars: carService,
		tmpl: mustLoadTemplates(),
	}

	// Cars pages
	mux.HandleFunc("/ui/cars", h.carsList)     // GET
	mux.HandleFunc("/ui/cars/new", h.carsNew)  // GET + POST
	mux.HandleFunc("/ui/cars/", h.carsActions) // POST actions: delete/reserve

	// Orders/Auth pages (placeholders)
	mux.HandleFunc("/ui/orders", h.ordersList) // GET
	mux.HandleFunc("/ui/login", h.login)       // GET
	mux.HandleFunc("/ui/register", h.register) // GET
}

func mustLoadTemplates() *template.Template {
	t := template.New("")

	patterns := []string{
		filepath.Join("web", "templates", "*.html"),
		filepath.Join("web", "templates", "cars", "*.html"),
		filepath.Join("web", "templates", "orders", "*.html"),
		filepath.Join("web", "templates", "auth", "*.html"),
	}

	var err error
	for _, p := range patterns {
		t, err = t.ParseGlob(p)
		if err != nil {
			panic(err)
		}
	}
	return t
}

// -------------------- Cars --------------------

type CarsListView struct {
	Title string
	Cars  []cars.Car
}

func (h *Handler) carsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	view := CarsListView{
		Title: "Cars",
		Cars:  h.cars.List(),
	}

	h.render(w, "cars_list.html", view)
}

type CarsNewView struct {
	Title string
	Error string
}

func (h *Handler) carsNew(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.render(w, "cars_new.html", CarsNewView{Title: "Add Car"})
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.render(w, "cars_new.html", CarsNewView{Title: "Add Car", Error: "Invalid form"})
			return
		}

		year, _ := strconv.Atoi(r.FormValue("year"))
		price, _ := strconv.Atoi(r.FormValue("price"))
		mileage, _ := strconv.Atoi(r.FormValue("mileage"))

		_, err := h.cars.Create(cars.CreateCarRequest{
			Brand:   r.FormValue("brand"),
			Model:   r.FormValue("model"),
			Year:    year,
			Price:   price,
			Mileage: mileage,
		})
		if err != nil {
			h.render(w, "cars_new.html", CarsNewView{Title: "Add Car", Error: "Validation error. Check fields."})
			return
		}

		http.Redirect(w, r, "/ui/cars", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// POST /ui/cars/{id}/delete
// POST /ui/cars/{id}/reserve
func (h *Handler) carsActions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/ui/cars/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	action := parts[1]
	switch action {
	case "delete":
		_ = h.cars.Delete(id)
		http.Redirect(w, r, "/ui/cars", http.StatusSeeOther)
		return

	case "reserve":
		status := cars.StatusReserved
		_, _ = h.cars.Update(id, cars.UpdateCarRequest{
			Status: &status,
		})
		http.Redirect(w, r, "/ui/cars", http.StatusSeeOther)
		return

	default:
		http.NotFound(w, r)
		return
	}
}

// -------------------- Orders/Auth placeholders --------------------

type SimplePage struct {
	Title string
	Note  string
}

func (h *Handler) ordersList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.render(w, "orders_list.html", SimplePage{
		Title: "Orders",
		Note:  "Orders UI will be implemented by Nurbol.",
	})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.render(w, "auth_login.html", SimplePage{
		Title: "Login",
		Note:  "Auth UI will be implemented by Ehson.",
	})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.render(w, "auth_register.html", SimplePage{
		Title: "Register",
		Note:  "Auth UI will be implemented by Ehson.",
	})
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
