package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"video-platform/config"
	"video-platform/index"
)

func ListVideos(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	videos := videoService.ListVideos()

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
	})

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
