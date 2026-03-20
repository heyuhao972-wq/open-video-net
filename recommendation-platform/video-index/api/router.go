package api

import (
	"video-index/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/video", handler.UpsertVideo)
	r.GET("/video/:id", handler.GetVideo)
	r.GET("/videos", handler.ListVideos)
	r.GET("/search", handler.Search)
	r.POST("/video/views", handler.UpdateViews)
}
