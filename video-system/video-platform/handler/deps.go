package handler

import (
	"video-platform/service"
	"video-platform/storage"
)

var (
	videoService  *service.VideoService
	uploadService *service.UploadService
	userService   *service.UserService
	storageClient *storage.StorageClient
)

func InitServices(vs *service.VideoService, us *service.UploadService, usvc *service.UserService, sc *storage.StorageClient) {
	videoService = vs
	uploadService = us
	userService = usvc
	storageClient = sc
}
