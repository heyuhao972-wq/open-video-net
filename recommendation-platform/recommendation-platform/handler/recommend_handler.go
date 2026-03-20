package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"recommendation-platform/config"
	"recommendation-platform/model"
	"recommendation-platform/pipeline"
	"recommendation-platform/repository"
	"recommendation-platform/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Recommend(c *gin.Context) {

	user := c.Query("user_id")
	if user == "" {
		user = c.Query("user")
	}
	recType := strings.TrimSpace(c.Query("type"))

	page := 1
	limit := 10
	if v := strings.TrimSpace(c.Query("page")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := strings.TrimSpace(c.Query("limit")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > 50 {
		limit = 50
	}
	fetchLimit := page * limit
	if fetchLimit > 100 {
		fetchLimit = 100
	}

	s := service.NewRecommendService()

	result := s.Recommend(user, recType, fetchLimit)
	start := (page - 1) * limit
	if start >= len(result) {
		c.JSON(200, gin.H{"videos": []string{}})
		return
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}
	result = result[start:end]

	out := make([]string, 0, len(result))
	for _, v := range result {
		if v.PlatformID == "" {
			continue
		}
		out = append(out, "video://"+v.PlatformID+"/"+v.ID)
	}

	c.JSON(200, gin.H{
		"videos": out,
	})

}

func getUserIDFromAuth(c *gin.Context) (string, bool) {
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return "", false
	}
	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	cfg := config.LoadConfig()

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return "", false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok && sub != "" {
			return sub, true
		}
	}
	return "", false
}

func AddBehavior(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	var b model.Behavior

	err := c.BindJSON(&b)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})

		return
	}

	b.UserID = userID
	if b.Timestamp == 0 {
		b.Timestamp = time.Now().Unix()
	}

	repository.AddBehavior(b)
	notifyFromBehavior(b)

	c.JSON(200, gin.H{
		"status": "ok",
	})

}

func FollowAuthor(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	var f model.Follow
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	f.UserID = userID
	if f.UserID == "" || f.AuthorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and author_id required",
		})
		return
	}

	repository.AddFollow(f)
	repository.AddNotification(model.Notification{
		UserID:  f.AuthorID,
		ActorID: f.UserID,
		Type:    "follow",
	})
	c.JSON(200, gin.H{"status": "ok"})
}

func UnfollowAuthor(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	var f model.Follow
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	f.UserID = userID
	if f.UserID == "" || f.AuthorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and author_id required",
		})
		return
	}

	repository.RemoveFollow(f.UserID, f.AuthorID)
	c.JSON(200, gin.H{"status": "ok"})
}

func GetMyLikes(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	likes := repository.GetLikesByUser(userID)
	videoIDs := make([]string, 0, len(likes))
	for _, b := range likes {
		if b.VideoID != "" {
			videoIDs = append(videoIDs, b.VideoID)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"videos": videoIDs,
	})
}

func GetMyFollows(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	follows := repository.GetFollowsByUser(userID)
	users := make([]string, 0, len(follows))
	for _, f := range follows {
		if f.AuthorID != "" {
			users = append(users, f.AuthorID)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetMyFollowers(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	follows := repository.GetFollowersByUser(userID)
	users := make([]string, 0, len(follows))
	for _, f := range follows {
		if f.UserID != "" {
			users = append(users, f.UserID)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetNotifications(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
	out := repository.GetNotificationsByUser(userID)
	c.JSON(http.StatusOK, gin.H{"notifications": out})
}

func MarkNotificationsRead(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
	var req struct {
		ID  int  `json:"id"`
		All bool `json:"all"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	if req.All {
		repository.MarkAllRead(userID)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	if req.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	repository.MarkNotificationRead(userID, req.ID)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func notifyFromBehavior(b model.Behavior) {
	if b.VideoID == "" || b.UserID == "" {
		return
	}
	authorID := resolveVideoAuthor(b.VideoID)
	if authorID == "" || authorID == b.UserID {
		return
	}
	switch b.Type {
	case "like", "share", "comment":
		repository.AddNotification(model.Notification{
			UserID:  authorID,
			ActorID: b.UserID,
			VideoID: b.VideoID,
			Type:    b.Type,
		})
	}
}

func resolveVideoAuthor(videoID string) string {
	videos := pipeline.LoadIndexVideos(pipeline.Query{})
	for _, v := range videos {
		if v.ID == videoID {
			return v.AuthorID
		}
	}
	return ""
}

func AddFavorite(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
	var req struct {
		VideoID string `json:"video_id"`
	}
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.VideoID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "video_id required",
		})
		return
	}
	repository.AddFavorite(model.Favorite{UserID: userID, VideoID: req.VideoID})
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func RemoveFavorite(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
	var req struct {
		VideoID string `json:"video_id"`
	}
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.VideoID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "video_id required",
		})
		return
	}
	repository.RemoveFavorite(userID, req.VideoID)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetMyFavorites(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
	favs := repository.GetFavoritesByUser(userID)
	videoIDs := make([]string, 0, len(favs))
	for _, f := range favs {
		if f.VideoID != "" {
			videoIDs = append(videoIDs, f.VideoID)
		}
	}
	c.JSON(http.StatusOK, gin.H{"videos": videoIDs})
}

func GetMyHistory(c *gin.Context) {
	userID, ok := getUserIDFromAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
	limit := 50
	history := repository.GetWatchHistoryByUser(userID, limit)

	platformMap := map[string]string{}
	for _, v := range pipeline.LoadIndexVideos(pipeline.Query{}) {
		if v.ID != "" {
			platformMap[v.ID] = v.PlatformID
		}
	}

	type item struct {
		VideoID   string `json:"video_id"`
		PlatformID string `json:"platform_id"`
		Timestamp int64  `json:"timestamp"`
	}
	out := make([]item, 0, len(history))
	for _, h := range history {
		out = append(out, item{
			VideoID:   h.VideoID,
			PlatformID: platformMap[h.VideoID],
			Timestamp: h.Timestamp,
		})
	}
	c.JSON(http.StatusOK, gin.H{"history": out})
}

func GetVideoStats(c *gin.Context) {
	videoID := strings.TrimSpace(c.Param("id"))
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "video_id required",
		})
		return
	}
	behaviors := repository.GetBehaviors()
	likes := 0
	shares := 0
	watches := 0
	notInterested := 0
	for _, b := range behaviors {
		if b.VideoID != videoID {
			continue
		}
		switch b.Type {
		case "like":
			likes++
		case "share":
			shares++
		case "watch":
			watches++
		case "not_interested":
			notInterested++
		}
	}
	favorites := repository.GetFavoriteCount(videoID)
	c.JSON(http.StatusOK, gin.H{
		"likes":          likes,
		"shares":         shares,
		"watches":        watches,
		"not_interested": notInterested,
		"favorites":      favorites,
	})
}
