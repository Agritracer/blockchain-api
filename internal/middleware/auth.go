package middleware

import (
	"agritrace-api/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

func Auth() gin.HandlerFunc {
	config.LoadConfig()

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		totpCode := c.GetHeader("X-TOTP-Code")

		// Kiểm tra API key
		if apiKey != config.Cfg.APIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid API Key"})
			c.Abort()
			return
		}

		// Kiểm tra mã TOTP
		if !totp.Validate(totpCode, config.Cfg.TOTPSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid TOTP Code"})
			c.Abort()
			return
		}

		// Nếu hợp lệ, tiếp tục xử lý request
		c.Next()
	}
}
