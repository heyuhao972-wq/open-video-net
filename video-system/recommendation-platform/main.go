package main

import (
	"fmt"

	"recommendation-platform/api"
	"recommendation-platform/config"
	"recommendation-platform/db"
	"recommendation-platform/middleware"
	"recommendation-platform/repository"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	r := gin.Default()
	r.Use(middleware.CORS())

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		panic(err)
	}
	repository.Init(database)

	api.RegisterRoutes(r)

	fmt.Println("recommendation service running:", cfg.Port)

	r.Run(":" + cfg.Port)

}
