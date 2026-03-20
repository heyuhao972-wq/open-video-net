package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"video-platform/config"

	"github.com/gin-gonic/gin"
)

func UploadVideo(c *gin.Context) {

	title := c.PostForm("title")
	description := c.PostForm("description")
	tagsRaw := c.PostForm("tags")
	if strings.TrimSpace(title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "title required",
		})
		return
	}

	if uploadService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "upload service not initialized",
		})
		return
	}

	file, err := c.FormFile("file")

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file required",
		})

		return
	}

	savePath := "./uploads"

	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "create upload dir failed",
		})
		return
	}

	safeName := filepath.Base(file.Filename)
	safeName = strings.ReplaceAll(safeName, "..", "")
	if safeName == "." || safeName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid filename",
		})
		return
	}

	path := filepath.Join(savePath, safeName)

	err = c.SaveUploadedFile(file, path)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "save failed",
		})

		return
	}

	coverPath := ""
	cover, _ := c.FormFile("cover")
	if cover != nil {
		coverDir := "./uploads/covers"
		if err := os.MkdirAll(coverDir, os.ModePerm); err == nil {
			coverName := filepath.Base(cover.Filename)
			coverName = strings.ReplaceAll(coverName, "..", "")
			if coverName != "" && coverName != "." {
				coverPath = filepath.Join(coverDir, coverName)
				_ = c.SaveUploadedFile(cover, coverPath)
			}
		}
	}

	tags := parseTags(tagsRaw)
	authorSignature := strings.TrimSpace(c.PostForm("author_signature"))
	videoHash := strings.TrimSpace(c.PostForm("video_hash"))
	authorTimestampRaw := strings.TrimSpace(c.PostForm("author_timestamp"))
	if authorSignature == "" || videoHash == "" || authorTimestampRaw == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "author_signature, video_hash, author_timestamp required",
		})
		return
	}
	authorTimestamp, err := strconv.ParseInt(authorTimestampRaw, 10, 64)
	if err != nil || authorTimestamp <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid author_timestamp",
		})
		return
	}
	cfg := config.LoadConfig()
	if len(cfg.AcceptTags) > 0 && !matchTags(tags, cfg.AcceptTags) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tags not allowed by platform policy",
		})
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user context",
		})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user context",
		})
		return
	}

	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}
	user, ok := userService.GetByID(userIDStr)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}

	video, err := uploadService.UploadVideo(title, description, tags, path, safeName, coverPath, user.ID, user.PublicKey, authorSignature, authorTimestamp, videoHash, cfg.PlatformID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"video": video,
	})

}

func parseTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func matchTags(input []string, allowed []string) bool {
	if len(input) == 0 || len(allowed) == 0 {
		return false
	}
	allow := map[string]bool{}
	for _, t := range allowed {
		allow[strings.ToLower(t)] = true
	}
	for _, t := range input {
		if allow[strings.ToLower(t)] {
			return true
		}
	}
	return false
}
