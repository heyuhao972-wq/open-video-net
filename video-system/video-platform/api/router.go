package api

import (
	"video-platform/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/health", handler.Health)

	r.POST("/register", handler.Register)
	r.POST("/login/challenge", handler.LoginChallenge)
	r.POST("/login", handler.Login)

	r.POST("/upload", handler.RequireAuth(), handler.UploadVideo)

	r.GET("/videos", handler.ListVideos)

	r.GET("/video/:id", handler.GetVideo)
	r.GET("/video/:id/stream", handler.StreamVideo)
	r.GET("/video/:id/manifest", handler.GetManifest)
	r.GET("/chunk/:hash", handler.GetChunk)

}
