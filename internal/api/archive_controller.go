package api

import (
	"archive_creator/internal/service"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ArchiveController struct {
	service *service.ArchiveService
}

func NewArchiveController(service *service.ArchiveService) *ArchiveController {
	return &ArchiveController{
		service: service,
	}
}

func (a *ArchiveController) CreateArchive(c *gin.Context) {
	id, err := a.service.CreateArchive()
	if err != nil {
		c.JSON(err.HttpStatusCode(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, CreateArchiveResponse{id})
}

func (a *ArchiveController) UpdateArchive(c *gin.Context) {

	archiveId := c.Param("id")

	var req UpdateArchiveDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.service.AddUrlToArchive(archiveId, req.FileUrl); err != nil {
		c.JSON(err.HttpStatusCode(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UpdateArchiveResponse{archiveId})
}

func (a *ArchiveController) GetArchiveStatus(c *gin.Context) {
	archiveId := c.Param("id")

	amount, status, url, err := a.service.GetArchiveStatus(archiveId)
	if err != nil {
		c.JSON(err.HttpStatusCode(), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetArchiveStatusResponse{status, amount, url})
}

func (a *ArchiveController) DownLoadArchive(c *gin.Context) {
	archiveId := c.Param("id")

	archivePath, err := a.service.GetArchivePath(archiveId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(archivePath))
	c.Header("Content-Type", "application/zip")
	c.File(archivePath)
}
