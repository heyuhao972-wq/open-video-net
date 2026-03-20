package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"video-platform/config"
	"video-platform/index"
	"video-platform/model"
)

func ListVideos(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	page, limit := parsePageLimit(c.Query("page"), c.Query("limit"))
	videos := videoService.ListVideos()
	videos = paginateVideos(videos, page, limit)

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
	})

}

func SearchVideos(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	q := strings.TrimSpace(c.Query("q"))
	tag := strings.TrimSpace(c.Query("tag"))
	page, limit := parsePageLimit(c.Query("page"), c.Query("limit"))
	videos := videoService.Search(q, tag)
	videos = paginateVideos(videos, page, limit)
	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
	})
}

func ListMyVideos(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
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

	page, limit := parsePageLimit(c.Query("page"), c.Query("limit"))
	videos := videoService.ListByAuthor(userIDStr)
	videos = paginateVideos(videos, page, limit)
	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
	})
}

func ListUserVideos(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}
	authorID := strings.TrimSpace(c.Param("id"))
	if authorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id required",
		})
		return
	}
	page, limit := parsePageLimit(c.Query("page"), c.Query("limit"))
	videos := videoService.ListByAuthor(authorID)
	videos = paginateVideos(videos, page, limit)
	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

func UpdateVideo(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "video service not initialized"})
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

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	video, ok := videoService.GetVideo(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	if video.AuthorID != userIDStr {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	tags := parseTags(req.Tags)
	updated, ok := videoService.UpdateVideo(id, strings.TrimSpace(req.Title), strings.TrimSpace(req.Description), tags)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"video": updated})
}

func DeleteMyVideo(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "video service not initialized"})
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

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	video, ok := videoService.GetVideo(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	if video.AuthorID != userIDStr {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	deleted, ok := videoService.DeleteVideo(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	if deleted.FilePath != "" {
		_ = os.Remove(deleted.FilePath)
	}
	if deleted.Manifest != "" {
		_ = os.Remove(deleted.Manifest)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func parsePageLimit(pageRaw string, limitRaw string) (int, int) {
	page := 1
	limit := 20
	if p, err := strconv.Atoi(pageRaw); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitRaw); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	return page, limit
}

func paginateVideos(videos []model.Video, page int, limit int) []model.Video {
	if limit <= 0 {
		return videos
	}
	start := (page - 1) * limit
	if start >= len(videos) {
		return []model.Video{}
	}
	end := start + limit
	if end > len(videos) {
		end = len(videos)
	}
	return videos[start:end]
}

func GetVideo(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	id := c.Param("id")
	video, ok := videoService.GetVideo(id)
	if ok {
		c.JSON(http.StatusOK, video)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "video not found",
	})

}

func StreamVideo(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	id := c.Param("id")
	video, ok := videoService.GetVideo(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "video not found",
		})
		return
	}

	if video.FilePath == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "video file missing",
		})
		return
	}

	c.Header("Content-Type", "video/mp4")
	c.File(video.FilePath)

	go func(id string) {
		cfg := config.LoadConfig()
		client := index.NewClient(cfg.IndexBase)
		_ = client.IncrementViews(id)
	}(video.ID)
}

func GetManifest(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	id := c.Param("id")
	video, ok := videoService.GetVideo(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "video not found",
		})
		return
	}

	if video.Manifest == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "manifest not found",
		})
		return
	}

	c.File(video.Manifest)
}

func GetChunk(c *gin.Context) {
	if storageClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "storage client not initialized",
		})
		return
	}

	hash := c.Param("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "hash required",
		})
		return
	}

	data, err := storageClient.GetChunk(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "chunk not found",
		})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}
