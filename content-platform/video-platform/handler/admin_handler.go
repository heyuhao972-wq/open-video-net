package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func BanUser(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user service not initialized"})
		return
	}
	var req struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.UserID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	_ = userService.BanUser(strings.TrimSpace(req.UserID), strings.TrimSpace(req.Reason))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func UnbanUser(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user service not initialized"})
		return
	}
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.UserID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	_ = userService.UnbanUser(strings.TrimSpace(req.UserID))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ListBans(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user service not initialized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bans": userService.ListBans()})
}

func AdminDeleteComment(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "comment service not initialized"})
		return
	}
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := commentService.Delete(id, "admin", true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func AdminDeleteVideo(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "video service not initialized"})
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	video, ok := videoService.DeleteVideo(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	if video.FilePath != "" {
		_ = os.Remove(video.FilePath)
	}
	if video.Manifest != "" {
		_ = os.Remove(video.Manifest)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ListReports(c *gin.Context) {
	if reportRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "report repo not initialized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reportRepo.List()})
}

func ListPendingVideos(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "video service not initialized"})
		return
	}
	status := strings.TrimSpace(c.Query("status"))
	if status == "" {
		status = "pending"
	}
	videos := videoService.ListVideosByStatus(status)
	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

func ReviewVideo(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "video service not initialized"})
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	var req struct {
		Status   string `json:"status"`
		Reason   string `json:"reason"`
		Reviewer string `json:"reviewer"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "approved"
	}
	video, ok := videoService.ReviewVideo(id, status, strings.TrimSpace(req.Reason), strings.TrimSpace(req.Reviewer))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"video": video})
}

func ListPendingComments(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "comment service not initialized"})
		return
	}
	status := strings.TrimSpace(c.Query("status"))
	if status == "" {
		status = "pending"
	}
	comments, err := commentService.ListByStatus(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func ReviewComment(c *gin.Context) {
	if commentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "comment service not initialized"})
		return
	}
	idStr := strings.TrimSpace(c.Param("id"))
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Status   string `json:"status"`
		Reason   string `json:"reason"`
		Reviewer string `json:"reviewer"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "approved"
	}
	comment, ok, err := commentService.Review(id, status, strings.TrimSpace(req.Reason), strings.TrimSpace(req.Reviewer))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comment": comment})
}
