package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = func() []byte {
	// Prefer env secret (recommended)
	if v := os.Getenv("JWT_SECRET"); v != "" {
		return []byte(v)
	}
	// fallback for local dev
	return []byte("my_secret_key_2026")
}()

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	usersMu sync.RWMutex
	usersDB = make(map[string]string) // username -> bcrypt hash
)

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

func GenerateJWT(username string) (string, error) {
	if username == "" {
		return "", errors.New("empty username")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString(jwtKey)
}

func ValidateToken(signedToken string) (string, error) {
	if signedToken == "" {
		return "", errors.New("empty token")
	}

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		// Validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	u, ok := claims["username"].(string)
	if !ok || u == "" {
		return "", fmt.Errorf("invalid token claims")
	}

	return u, nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	hashed, err := HashPassword(input.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	usersMu.Lock()
	_, exists := usersDB[input.Username]
	if exists {
		usersMu.Unlock()
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}
	usersDB[input.Username] = hashed
	usersMu.Unlock()

	// Background log (goroutine requirement)
	go func(name string) {
		time.Sleep(1 * time.Second)
		fmt.Printf("[LOG] User %s registered in background\n", name)
	}(input.Username)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "registration successful"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	usersMu.RLock()
	storedHash, exists := usersDB[input.Username]
	usersMu.RUnlock()

	if !exists || !CheckPasswordHash(input.Password, storedHash) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(input.Username)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
