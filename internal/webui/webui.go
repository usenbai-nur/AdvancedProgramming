package webui

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"AdvancedProgramming/internal/auth"
	"AdvancedProgramming/internal/cars"
)

type Handler struct {
	cars *cars.Service
	tmpl *template.Template
}

type BaseView struct {
	Title string
}

func Register(mux *http.ServeMux, carService *cars.Service) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	h := &Handler{cars: carService, tmpl: mustLoadTemplates()}

	mux.HandleFunc("/ui/cars", h.carsList)
	mux.HandleFunc("/ui/cars/new", h.carsNew)
	mux.HandleFunc("/ui/cars/", h.carsActions)
	mux.HandleFunc("/ui/orders", h.ordersList)
	mux.HandleFunc("/ui/login", h.login)
	mux.HandleFunc("/ui/register", h.register)
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

type CarsListView struct {
	BaseView
	Cars []cars.Car
}

func (h *Handler) carsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.render(w, "cars_list.html", CarsListView{BaseView: BaseView{Title: "Cars"}, Cars: h.cars.List()})
}

type CarsNewView struct {
	BaseView
	Error string
}

func (h *Handler) carsNew(w http.ResponseWriter, r *http.Request) {
	view := CarsNewView{BaseView: BaseView{Title: "Add Car"}}
	switch r.Method {
	case http.MethodGet:
		h.render(w, "cars_new.html", view)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			view.Error = "Invalid form"
			h.render(w, "cars_new.html", view)
			return
		}
		year, _ := strconv.Atoi(r.FormValue("year"))
		price, _ := strconv.Atoi(r.FormValue("price"))
		mileage, _ := strconv.Atoi(r.FormValue("mileage"))
		_, err := h.cars.Create(cars.CreateCarRequest{
			Brand: r.FormValue("brand"), Model: r.FormValue("model"), Year: year, Price: price, Mileage: mileage,
		})
		if err != nil {
			view.Error = "Validation error. Check fields."
			h.render(w, "cars_new.html", view)
			return
		}
		http.Redirect(w, r, "/ui/cars", http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

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
	switch parts[1] {
	case "delete":
		_ = h.cars.Delete(id)
	case "reserve":
		status := cars.StatusReserved
		_, _ = h.cars.Update(id, cars.UpdateCarRequest{Status: &status})
	default:
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui/cars", http.StatusSeeOther)
}

type SimplePage struct {
	BaseView
	Note string
}

type AuthPage struct {
	BaseView
	Error   string
	Success string
	Token   string
}

func (h *Handler) ordersList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.render(w, "orders_list.html", SimplePage{BaseView: BaseView{Title: "Orders"}, Note: "Orders are managed via API. Admin can view and process all orders."})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	view := AuthPage{BaseView: BaseView{Title: "Login"}}
	if r.Method == http.MethodGet {
		h.render(w, "auth_login.html", view)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		view.Error = "Invalid form"
		h.render(w, "auth_login.html", view)
		return
	}
	token, user, err := auth.LoginUser(auth.LoginRequest{Username: r.FormValue("username"), Password: r.FormValue("password")})
	if err != nil {
		view.Error = err.Error()
		h.render(w, "auth_login.html", view)
		return
	}
	view.Success = "Welcome, " + user.Username + " (" + string(user.Role) + ")"
	view.Token = token
	h.render(w, "auth_login.html", view)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	view := AuthPage{BaseView: BaseView{Title: "Register"}}
	if r.Method == http.MethodGet {
		h.render(w, "auth_register.html", view)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		view.Error = "Invalid form"
		h.render(w, "auth_register.html", view)
		return
	}
	role := "user"
	if r.FormValue("role") == "admin" {
		role = "admin"
	}
	_, err := auth.RegisterUser(auth.RegisterRequest{Username: r.FormValue("username"), Password: r.FormValue("password"), Role: role, AdminKey: r.FormValue("admin_key")})
	if err != nil {
		view.Error = err.Error()
		h.render(w, "auth_register.html", view)
		return
	}
	view.Success = "Registration successful. You can login now."
	h.render(w, "auth_register.html", view)
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
