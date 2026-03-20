package repository

import "recommendation-platform/model"

var behaviors []model.Behavior
var follows []model.Follow
var favorites []model.Favorite

func AddBehavior(b model.Behavior) {

	behaviors = append(behaviors, b)

}

func GetBehaviors() []model.Behavior {

	return behaviors

}

func GetBehaviorsByUser(userID string) []model.Behavior {
	out := []model.Behavior{}
	for _, b := range behaviors {
		if b.UserID == userID {
			out = append(out, b)
		}
	}
	return out
}

func GetLikesByUser(userID string) []model.Behavior {
	out := []model.Behavior{}
	for _, b := range behaviors {
		if b.UserID == userID && b.Type == "like" {
			out = append(out, b)
		}
	}
	return out
}

func GetWatchHistoryByUser(userID string, limit int) []model.Behavior {
	out := []model.Behavior{}
	for i := len(behaviors) - 1; i >= 0; i-- {
		b := behaviors[i]
		if b.UserID == userID && b.Type == "watch" {
			out = append(out, b)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out
}

func AddFollow(f model.Follow) {

	if !f.Active {
		f.Active = true
	}
	follows = append(follows, f)

}

func GetFollows() []model.Follow {

	return follows

}

func GetFollowsByUser(userID string) []model.Follow {
	out := []model.Follow{}
	for _, f := range follows {
		if f.UserID == userID {
			out = append(out, f)
		}
	}
	return out
}

func GetFollowersByUser(authorID string) []model.Follow {
	out := []model.Follow{}
	for _, f := range follows {
		if f.AuthorID == authorID {
			out = append(out, f)
		}
	}
	return out
}

func AddFavorite(f model.Favorite) {
	for _, existing := range favorites {
		if existing.UserID == f.UserID && existing.VideoID == f.VideoID {
			return
		}
	}
	favorites = append(favorites, f)
}

func RemoveFavorite(userID string, videoID string) {
	for i := len(favorites) - 1; i >= 0; i-- {
		if favorites[i].UserID == userID && favorites[i].VideoID == videoID {
			favorites = append(favorites[:i], favorites[i+1:]...)
			return
		}
	}
}

func GetFavoritesByUser(userID string) []model.Favorite {
	out := []model.Favorite{}
	for _, f := range favorites {
		if f.UserID == userID {
			out = append(out, f)
		}
	}
	return out
}

func GetFavoriteCount(videoID string) int {
	count := 0
	for _, f := range favorites {
		if f.VideoID == videoID {
			count++
		}
	}
	return count
}

func RemoveFollow(userID string, authorID string) {
	for i := len(follows) - 1; i >= 0; i-- {
		if follows[i].UserID == userID && follows[i].AuthorID == authorID {
			follows = append(follows[:i], follows[i+1:]...)
			return
		}
	}
}
