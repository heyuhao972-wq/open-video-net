package main

import (
	"fmt"

	"streaming-service/api"
	"streaming-service/config"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	r := gin.Default()

	api.RegisterRoutes(r)

	fmt.Println("Streaming service running on port:", cfg.Port)

	err := r.Run(":" + cfg.Port)

	if err != nil {
		panic(err)
	}

}
