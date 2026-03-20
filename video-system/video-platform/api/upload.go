package api

import (
	"video-platform/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(r *gin.Engine) {

	r.POST("/upload", handler.UploadVideo)

}
