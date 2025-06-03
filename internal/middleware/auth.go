package middleware

import (
	"agritrace-api/internal/config"
	"net/http"

	"github.com/pquerna/otp/totp"
)

func Auth(next http.Handler) http.Handler {
	config.LoadConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		totpCode := r.Header.Get("X-TOTP-Code")

		// Kiểm tra API key
		if apiKey != config.Cfg.APIKey {
			http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
			return
		}

		// Kiểm tra mã TOTP
		if !totp.Validate(totpCode, config.Cfg.TOTPSecret) {
			http.Error(w, "Unauthorized: Invalid TOTP Code", http.StatusUnauthorized)
			return
		}

		// Nếu hợp lệ, tiếp tục xử lý request
		next.ServeHTTP(w, r)
	})
}
