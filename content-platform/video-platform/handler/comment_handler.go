package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"video-platform/model"
)

type commentRequest struct {
	VideoID  string `json:"video_id"`
	Content  string `json:"content"`
	ParentID int    `json:"parent_id"`
}

func CreateComment(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "comment service not initialized",
		})
		return
	}

	var req commentRequest
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

	comment, err := commentService.Create(strings.TrimSpace(req.VideoID), userIDStr, req.Content, req.ParentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     comment.ID,
		"status": comment.Status,
	})
}

func GetVideoComments(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "comment service not initialized",
		})
		return
	}

	videoID := strings.TrimSpace(c.Param("id"))
	page, limit := parsePageLimit(c.Query("page"), c.Query("limit"))
	comments, err := commentService.List(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, paginateComments(comments, page, limit))
}

func GetCommentCount(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "comment service not initialized",
		})
		return
	}

	videoID := strings.TrimSpace(c.Param("id"))
	count, err := commentService.Count(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func DeleteComment(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "comment service not initialized",
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

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if err := commentService.Delete(id, userIDStr, false); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func LikeComment(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "comment service not initialized",
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

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	_, _, err = commentService.Like(id, userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func paginateComments(comments []model.Comment, page int, limit int) []model.Comment {
	if limit <= 0 {
		return comments
	}
	start := (page - 1) * limit
	if start >= len(comments) {
		return []model.Comment{}
	}
	end := start + limit
	if end > len(comments) {
		end = len(comments)
	}
	return comments[start:end]
}
