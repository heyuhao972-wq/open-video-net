package handler

import (
	"net/http"
	"strings"

	"recommendation-platform/config"
	"recommendation-platform/model"
	"recommendation-platform/repository"
	"recommendation-platform/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Recommend(c *gin.Context) {

	user := c.Query("user_id")
	if user == "" {
		user = c.Query("user")
	}

	s := service.NewRecommendService()

	result := s.Recommend(user)

	out := make([]string, 0, len(result))
	for _, v := range result {
		if v.PlatformID == "" {
			continue
		}
		out = append(out, "video://"+v.PlatformID+"/"+v.ID)
	}

	c.JSON(200, gin.H{
		"videos": out,
	})

}

func AddBehavior(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
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
		return
	}

	var b model.Behavior

	err = c.BindJSON(&b)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})

		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok && sub != "" {
			b.UserID = sub
		}
	}

	repository.AddBehavior(b)

	c.JSON(200, gin.H{
		"status": "ok",
	})

}

func FollowAuthor(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
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
		return
	}

	var f model.Follow
	if err = c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok && sub != "" {
			f.UserID = sub
		}
	}
	if f.UserID == "" || f.AuthorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and author_id required",
		})
		return
	}

	repository.AddFollow(f)
	c.JSON(200, gin.H{"status": "ok"})
}

func UnfollowAuthor(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
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
		return
	}

	var f model.Follow
	if err = c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok && sub != "" {
			f.UserID = sub
		}
	}
	if f.UserID == "" || f.AuthorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and author_id required",
		})
		return
	}

	repository.RemoveFollow(f.UserID, f.AuthorID)
	c.JSON(200, gin.H{"status": "ok"})
}
