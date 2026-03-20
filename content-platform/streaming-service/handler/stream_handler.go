package handler

import (
	"io"
	"net/http"

	"streaming-service/config"
	"streaming-service/service"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {

	c.JSON(200, gin.H{
		"status": "ok",
	})

}

func StreamVideo(c *gin.Context) {

	videoId := c.Param("videoId")
	platformID := c.Query("platform")

	cfg := config.LoadConfig()
	platformBase := cfg.PlatformBase
	p2pBase := cfg.P2PBase
	if platformID != "" && cfg.PlatformMap != nil {
		if v, ok := cfg.PlatformMap[platformID]; ok {
			platformBase = v
		}
	}
	if platformID != "" && cfg.P2PMap != nil {
		if v, ok := cfg.P2PMap[platformID]; ok {
			p2pBase = v
		}
	}

	streamService := service.NewStreamService(platformBase, p2pBase, cfg.MaxParallel, cfg.CacheSize, cfg.ChunkTimeoutMs, cfg.ChunkRetry)

	reader, err := streamService.GetVideoStream(videoId)

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "video not found",
		})

		return
	}

	c.Header("Content-Type", "video/mp4")

	c.Stream(func(w io.Writer) bool {

		buf := make([]byte, 1024*64)

		n, readErr := reader.Read(buf)

		if n > 0 {
			if _, err := w.Write(buf[:n]); err != nil {
				return false
			}
			return true
		}

		if readErr == io.EOF {
			return false
		}
		if readErr != nil {
			return false
		}
		return false
	})

}
