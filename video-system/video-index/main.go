package main

import (
	"fmt"

	"video-index/api"
	"video-index/config"
	"video-index/handler"
	"video-index/repository"
	"video-index/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	repo := repository.NewVideoRepository()
	videoService := service.NewVideoService(repo)
	handler.InitServices(videoService)

	r := gin.Default()
	api.RegisterRoutes(r)

	fmt.Println("video-index running on port:", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
