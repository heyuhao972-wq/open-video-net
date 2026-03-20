package main

import (
	"fmt"

	"video-platform/api"
	"video-platform/config"
	"video-platform/db"
	"video-platform/handler"
	"video-platform/index"
	"video-platform/middleware"
	"video-platform/repository"
	"video-platform/service"
	"video-platform/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		panic(err)
	}

	repo := repository.NewVideoRepository(database)
	userRepo := repository.NewUserRepository(database)
	videoService := service.NewVideoService(repo)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	storageClient, err := storage.NewStorageClient("./data/video-storage", 1024*1024)
	if err != nil {
		panic(err)
	}
	var indexClient *index.Client
	if cfg.IndexBase != "" {
		indexClient = index.NewClient(cfg.IndexBase)
	}
	uploadService := service.NewUploadService(videoService, storageClient, indexClient)
	handler.InitServices(videoService, uploadService, userService, storageClient)

	r := gin.Default()
	r.Use(middleware.CORS())

	api.RegisterRoutes(r)

	fmt.Println("Video Platform running on port:", cfg.Port)

	err = r.Run(":" + cfg.Port)
	if err != nil {
		panic(err)
	}

}
