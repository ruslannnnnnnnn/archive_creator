package api

import (
	archiveStorage "archive_creator/internal/archive_storage"
	"archive_creator/internal/config"
	"archive_creator/internal/service"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, config *config.Config) {

	archiveController := NewArchiveController(
		service.NewArchiveService(
			archiveStorage.NewStorage(),
			config,
		),
	)

	r.GET("/api/archive/:id/status", archiveController.GetArchiveStatus)
	r.GET("/api/archive/:id/download", archiveController.DownLoadArchive)
	r.POST("/api/archive", archiveController.CreateArchive)
	r.PATCH("/api/archive/:id/add-link", archiveController.UpdateArchive)
}
