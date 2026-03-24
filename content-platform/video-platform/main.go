package main

import (
	"fmt"

	"video-platform/api"
	"video-platform/config"
	"video-platform/db"
	"video-platform/handler"
	"video-platform/index"
	"video-platform/repository"
	"video-platform/service"
	"video-platform/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	repo := repository.NewVideoRepository()
	userRepo := repository.NewUserRepository()
	reportRepo := repository.NewReportRepository()
	sqlDB, err := db.Open(cfg.DBPath)
	if err != nil {
		panic(err)
	}
	if err := db.InitCommentTables(sqlDB); err != nil {
		panic(err)
	}
	commentRepo := repository.NewCommentRepository(sqlDB)
	videoService := service.NewVideoService(repo)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	commentService := service.NewCommentService(commentRepo)
	var storageClient *storage.StorageClient
	if cfg.StorageBase == "" {
		storageClient, err = storage.NewStorageClient("./data/video-storage", 1024*1024)
		if err != nil {
			panic(err)
		}
	}
	var indexClient *index.Client
	if cfg.IndexBase != "" {
		indexClient = index.NewClient(cfg.IndexBase)
	}
	uploadService := service.NewUploadService(videoService, indexClient)
	handler.InitServices(videoService, uploadService, userService, storageClient, commentService, reportRepo)

	r := gin.Default()
	r.Use(api.CORSMiddleware())

	r.Static("/uploads", "./uploads")

	api.RegisterRoutes(r)

	fmt.Println("Video Platform running on port:", cfg.Port)

	err = r.Run(":" + cfg.Port)
	if err != nil {
		panic(err)
	}

}
