package api

import (
	"recommendation-platform/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/health", handler.Health)
	r.GET("/recommend", handler.Recommend)

	r.POST("/behavior", handler.AddBehavior)
	r.POST("/follow", handler.FollowAuthor)
	r.POST("/unfollow", handler.UnfollowAuthor)

}
