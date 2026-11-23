package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserRole string `json:"user_role,omitempty"`
	jwt.RegisteredClaims
}

func authMiddleware (jwtSecret []byte, expectedIssuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format"},
			)
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims := &AuthClaims{}

		token, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (any, error) {
				// Enforce HMAC signing method so attackers cannot switch algorithms (e.g., "none").
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return jwtSecret, nil
			},
			jwt.WithIssuer(expectedIssuer),
			jwt.WithLeeway(5*time.Second),
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if claims.Subject == "" {
			log.Println("token missing subject")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing subject"})
			return
		}
		if claims.UserRole == "" {
			log.Println("token missing role")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing role"})
			return
		}

		// Store verified user ID in Gin request context
		c.Set("user_id", claims.Subject)
		c.Set("user_role", claims.UserRole)
		c.Next()
	}
}