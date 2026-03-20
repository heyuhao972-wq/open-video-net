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
	r.POST("/comment", handler.RequireAuth(), handler.CreateComment)
	r.POST("/report/video", handler.RequireAuth(), handler.ReportVideo)
	r.POST("/report/comment", handler.RequireAuth(), handler.ReportComment)

	r.GET("/videos", handler.ListVideos)
	r.GET("/search", handler.SearchVideos)
	r.GET("/me/videos", handler.RequireAuth(), handler.ListMyVideos)
	r.GET("/user/:id/videos", handler.ListUserVideos)
	r.GET("/video/:id/comments", handler.GetVideoComments)
	r.GET("/video/:id/comments/count", handler.GetCommentCount)

	r.GET("/video/:id", handler.GetVideo)
	r.PUT("/video/:id", handler.RequireAuth(), handler.UpdateVideo)
	r.DELETE("/video/:id", handler.RequireAuth(), handler.DeleteMyVideo)
	r.GET("/video/:id/stream", handler.StreamVideo)
	r.GET("/video/:id/manifest", handler.GetManifest)
	r.GET("/chunk/:hash", handler.GetChunk)
	r.DELETE("/comment/:id", handler.RequireAuth(), handler.DeleteComment)
	r.POST("/comment/:id/like", handler.RequireAuth(), handler.LikeComment)

	r.GET("/profile/:id", handler.GetProfile)
	r.GET("/me/profile", handler.RequireAuth(), handler.GetMyProfile)
	r.POST("/profile", handler.RequireAuth(), handler.UpdateProfile)
	r.GET("/users/search", handler.SearchUsers)

	r.GET("/admin/reports", handler.RequireAdmin(), handler.ListReports)
	r.POST("/admin/ban", handler.RequireAdmin(), handler.BanUser)
	r.POST("/admin/unban", handler.RequireAdmin(), handler.UnbanUser)
	r.GET("/admin/bans", handler.RequireAdmin(), handler.ListBans)
	r.DELETE("/admin/comment/:id", handler.RequireAdmin(), handler.AdminDeleteComment)
	r.DELETE("/admin/video/:id", handler.RequireAdmin(), handler.AdminDeleteVideo)

}
