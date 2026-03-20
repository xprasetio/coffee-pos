package middleware

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xprasetio/coffee-pos/pkg/response"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Token tidak ditemukan")
			c.Abort()
			return
		}

		// 2. Cek prefix "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "Format token tidak valid")
			c.Abort()
			return
		}

		// 3. Potong prefix "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. Parse dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "Token tidak valid atau sudah expired")
			c.Abort()
			return
		}

		// 5. Ekstrak claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "Token tidak valid")
			c.Abort()
			return
		}

		// 6. Ambil user_id dan role dari claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			response.Unauthorized(c, "Token tidak valid")
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			response.Unauthorized(c, "Token tidak valid")
			c.Abort()
			return
		}

		// 7. Simpan ke Gin context
		c.Set("user_id", userID)
		c.Set("role", role)

		// 8. Lanjut ke handler berikutnya
		c.Next()
	}
}

// RoleMiddleware creates a role-based authorization middleware
// It checks if the user's role is in the list of allowed roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil role dari Gin context
		role := c.GetString("role")
		if role == "" {
			response.Unauthorized(c, "Token tidak valid")
			c.Abort()
			return
		}

		// 2. Cek apakah role ada di dalam allowedRoles
		if !slices.Contains(allowedRoles, role) {
			response.Forbidden(c, "Akses ditolak")
			c.Abort()
			return
		}

		// 3. Lanjut ke handler berikutnya
		c.Next()
	}
}
