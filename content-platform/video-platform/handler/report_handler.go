package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type reportRequest struct {
	VideoID   string `json:"video_id"`
	CommentID string `json:"comment_id"`
	Reason    string `json:"reason"`
}

func ReportVideo(c *gin.Context) {
	if reportRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "report repo not initialized"})
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return
	}
	userIDStr, _ := userID.(string)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return
	}
	var req reportRequest
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.VideoID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video_id required"})
		return
	}
	reportRepo.Add("video", strings.TrimSpace(req.VideoID), userIDStr, strings.TrimSpace(req.Reason))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ReportComment(c *gin.Context) {
	if reportRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "report repo not initialized"})
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return
	}
	userIDStr, _ := userID.(string)
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return
	}
	var req reportRequest
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.CommentID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comment_id required"})
		return
	}
	reportRepo.Add("comment", strings.TrimSpace(req.CommentID), userIDStr, strings.TrimSpace(req.Reason))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
