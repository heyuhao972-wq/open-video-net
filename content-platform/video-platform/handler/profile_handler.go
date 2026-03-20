package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type profileRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

func GetProfile(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id required",
		})
		return
	}

	user, ok := userService.GetByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func GetMyProfile(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user context",
		})
		return
	}
	userIDStr, _ := userID.(string)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user context",
		})
		return
	}

	user, ok := userService.GetByID(userIDStr)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func UpdateProfile(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}

	var req profileRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user context",
		})
		return
	}
	userIDStr, _ := userID.(string)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user context",
		})
		return
	}

	user, err := userService.UpdateProfile(
		userIDStr,
		strings.TrimSpace(req.Nickname),
		strings.TrimSpace(req.AvatarURL),
		strings.TrimSpace(req.Bio),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func SearchUsers(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}

	q := strings.TrimSpace(c.Query("q"))
	users := userService.SearchUsers(q)
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
