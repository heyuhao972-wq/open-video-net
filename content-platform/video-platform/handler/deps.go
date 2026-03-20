package handler

import (
	"video-platform/repository"
	"video-platform/service"
	"video-platform/storage"
)

var (
	videoService  *service.VideoService
	uploadService *service.UploadService
	userService   *service.UserService
	storageClient *storage.StorageClient
	commentService *service.CommentService
	reportRepo *repository.ReportRepository
)

func InitServices(vs *service.VideoService, us *service.UploadService, usvc *service.UserService, sc *storage.StorageClient, cs *service.CommentService, rr *repository.ReportRepository) {
	videoService = vs
	uploadService = us
	userService = usvc
	storageClient = sc
	commentService = cs
	reportRepo = rr
}
