package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Role      Role   `json:"role"`
	Favorites []int  `json:"favorites,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
	AdminKey string `json:"admin_key,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string
	Role     Role
}

var jwtKey = func() []byte {
	if v := os.Getenv("JWT_SECRET"); v != "" {
		return []byte(v)
	}
	return []byte("my_secret_key_2026")
}()

var (
	usersMu sync.RWMutex
	nextID  int64
	usersDB = make(map[string]UserRecord) // username -> record
)

type UserRecord struct {
	ID           int
	Username     string
	PasswordHash string
	Role         Role
	Favorites    []int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateJWT(username string, role Role) (string, error) {
	if username == "" {
		return "", errors.New("empty username")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString(jwtKey)
}

func ValidateToken(signedToken string) (Claims, error) {
	if signedToken == "" {
		return Claims{}, errors.New("empty token")
	}

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}

	u, ok := claims["username"].(string)
	if !ok || u == "" {
		return Claims{}, fmt.Errorf("invalid token claims")
	}

	roleStr, _ := claims["role"].(string)
	role := Role(strings.ToLower(strings.TrimSpace(roleStr)))
	if role != RoleAdmin {
		role = RoleUser
	}

	return Claims{Username: u, Role: role}, nil
}

func RegisterUser(req RegisterRequest) (User, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		return User{}, errors.New("username and password required")
	}

	role := RoleUser
	if strings.EqualFold(strings.TrimSpace(req.Role), string(RoleAdmin)) {
		adminSecret := os.Getenv("ADMIN_REGISTRATION_KEY")
		if adminSecret == "" || req.AdminKey != adminSecret {
			return User{}, errors.New("invalid admin key")
		}
		role = RoleAdmin
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return User{}, errors.New("failed to hash password")
	}

	now := time.Now().UTC()
	usersMu.Lock()
	defer usersMu.Unlock()
	if _, exists := usersDB[username]; exists {
		return User{}, errors.New("user already exists")
	}

	id := int(atomic.AddInt64(&nextID, 1))
	rec := UserRecord{
		ID:           id,
		Username:     username,
		PasswordHash: hashed,
		Role:         role,
		Favorites:    []int{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	usersDB[username] = rec

	go func(name string) {
		time.Sleep(1 * time.Second)
		fmt.Printf("[LOG] User %s registered in background\n", name)
	}(username)

	return toUser(rec), nil
}

func LoginUser(req LoginRequest) (string, User, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		return "", User{}, errors.New("username and password required")
	}

	usersMu.RLock()
	rec, exists := usersDB[username]
	usersMu.RUnlock()

	if !exists || !CheckPasswordHash(password, rec.PasswordHash) {
		return "", User{}, errors.New("invalid credentials")
	}

	token, err := GenerateJWT(rec.Username, rec.Role)
	if err != nil {
		return "", User{}, errors.New("failed to generate token")
	}

	return token, toUser(rec), nil
}

func GetUserByUsername(username string) (User, bool) {
	usersMu.RLock()
	rec, exists := usersDB[username]
	usersMu.RUnlock()
	if !exists {
		return User{}, false
	}
	return toUser(rec), true
}

func AddFavorite(username string, carID int) (User, error) {
	if carID <= 0 {
		return User{}, errors.New("invalid car id")
	}
	usersMu.Lock()
	defer usersMu.Unlock()
	rec, ok := usersDB[username]
	if !ok {
		return User{}, errors.New("user not found")
	}
	for _, id := range rec.Favorites {
		if id == carID {
			return toUser(rec), nil
		}
	}
	rec.Favorites = append(rec.Favorites, carID)
	rec.UpdatedAt = time.Now().UTC()
	usersDB[username] = rec
	return toUser(rec), nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	user, err := RegisterUser(req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "registration successful",
		"user":    user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, user, err := LoginUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"token": token, "user": user})
}

func toUser(rec UserRecord) User {
	return User{
		ID:        rec.ID,
		Username:  rec.Username,
		Role:      rec.Role,
		Favorites: append([]int(nil), rec.Favorites...),
		CreatedAt: rec.CreatedAt.Format(time.RFC3339),
		UpdatedAt: rec.UpdatedAt.Format(time.RFC3339),
	}
}
