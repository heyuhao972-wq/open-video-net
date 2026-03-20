package main

import (
	"fmt"

	"recommendation-platform/api"
	"recommendation-platform/config"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	r := gin.Default()
	r.Use(api.CORSMiddleware())

	api.RegisterRoutes(r)

	fmt.Println("recommendation service running:", cfg.Port)

	r.Run(":" + cfg.Port)

}
