package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"video-platform/config"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		cfg := config.LoadConfig()

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok && sub != "" {
				c.Set("user_id", sub)
				if userService != nil {
					if banned, reason := userService.IsBanned(sub); banned {
						c.JSON(http.StatusForbidden, gin.H{
							"error": "user banned",
							"reason": reason,
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.LoadConfig()
		token := c.GetHeader("X-Admin-Token")
		if cfg.AdminToken == "" || token == "" || token != cfg.AdminToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "admin token required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
