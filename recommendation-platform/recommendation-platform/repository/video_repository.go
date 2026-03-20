package repository

import "recommendation-platform/model"

var videos = []model.Video{
	{ID: "1", Title: "AI Introduction", Views: 100},
	{ID: "2", Title: "Machine Learning", Views: 200},
	{ID: "3", Title: "Go Tutorial", Views: 80},
	{ID: "4", Title: "Distributed Systems", Views: 120},
}

func GetVideos() []model.Video {

	return videos

}
