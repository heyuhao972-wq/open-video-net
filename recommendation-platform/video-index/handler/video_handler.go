package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"video-index/model"
	"video-index/service"
)

var videoService *service.VideoService

func InitServices(vs *service.VideoService) {
	videoService = vs
}

func UpsertVideo(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	var v model.Video
	if err := c.BindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	if v.ID == "" || v.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id and title required",
		})
		return
	}

	if v.Tags == nil {
		v.Tags = []string{}
	}

	videoService.Save(v)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
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
	v, ok := videoService.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "video not found",
		})
		return
	}

	c.JSON(http.StatusOK, v)
}

func ListVideos(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"videos": videoService.List(),
	})
}

func Search(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	q := c.Query("q")
	results := videoService.Search(q)
	c.JSON(http.StatusOK, gin.H{
		"videos": results,
	})
}

type viewsRequest struct {
	ID    string `json:"id"`
	Op    string `json:"op"`
	Field string `json:"field"`
	Views int    `json:"views"`
}

func UpdateViews(c *gin.Context) {
	if videoService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "video service not initialized",
		})
		return
	}

	var req viewsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	v, ok := videoService.Get(req.ID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "video not found",
		})
		return
	}

	if req.Op == "inc" && req.Field == "views" {
		v.Views += req.Views
		videoService.Save(v)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
