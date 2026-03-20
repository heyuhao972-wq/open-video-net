package api

import (
	"streaming-service/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/health", handler.Health)

	r.GET("/stream/:videoId", handler.StreamVideo)

}
