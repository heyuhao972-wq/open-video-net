package api

import (
	"recommendation-platform/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/health", handler.Health)
	r.OPTIONS("/*path", func(c *gin.Context) { c.Status(204) })

	r.GET("/recommend", handler.Recommend)

	r.POST("/behavior", handler.AddBehavior)
	r.POST("/follow", handler.FollowAuthor)
	r.POST("/unfollow", handler.UnfollowAuthor)
	r.GET("/me/likes", handler.GetMyLikes)
	r.GET("/me/follows", handler.GetMyFollows)
	r.GET("/me/followers", handler.GetMyFollowers)
	r.POST("/favorite", handler.AddFavorite)
	r.POST("/unfavorite", handler.RemoveFavorite)
	r.GET("/me/favorites", handler.GetMyFavorites)
	r.GET("/me/history", handler.GetMyHistory)
	r.GET("/video/:id/stats", handler.GetVideoStats)
	r.GET("/notifications", handler.GetNotifications)
	r.POST("/notifications/read", handler.MarkNotificationsRead)

}
