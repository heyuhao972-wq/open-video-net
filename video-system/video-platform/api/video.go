package api

import (
	"video-platform/handler"

	"github.com/gin-gonic/gin"
)

func RegisterVideoRoutes(r *gin.Engine) {

	r.GET("/videos", handler.ListVideos)

	r.GET("/video/:id", handler.GetVideo)

}
