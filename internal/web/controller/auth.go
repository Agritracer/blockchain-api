package controller

import (
	"agritrace/internal/config"
	"encoding/json"
	"net/http"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var users = map[string]string{
	"admin": "$2y$10$kL4rzwXNkk5NO2EicLOAcemAxKVTTNV4qLpNd.UGuOjwX/Hm/fsg2", // 123456
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	config.LoadConfig()

	jwtKey = []byte(config.Cfg.JWTToken)

	if r.Method != http.MethodPost {
		http.Error(w, "Phải POST mới login được", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	storedHash, ok := users[creds.Username]
	if !ok {
		http.Error(w, "Sai username hoặc password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Sai username hoặc password", http.StatusUnauthorized)
		return
	}

	// Tạo token, trả về JSON
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Lỗi server", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Đăng nhập thành công",
	})
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	tplPath := filepath.Join(basepath, "../templates/login.html")

	t, err := template.ParseFiles(tplPath)
	if err != nil {
		http.Error(w, "Lỗi tải template", http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
