package service

import (
	"archive_creator/internal/api/api_error"
	archiveStorage "archive_creator/internal/archive_storage"
	fs "archive_creator/internal/archive_storage/filesystem"
	"archive_creator/internal/archive_storage/helpers"
	"archive_creator/internal/config"
	"fmt"
	"strconv"
)

const (
	TooMuchObjects       = "archive objects limit reached"
	ArchiveNotFound      = "archive not found"
	TooMuchArchives      = "server is busy"
	UnsupportedMediaType = "file format/mime not supported"
)

type IArchiveService interface {
	GetArchiveStatus(archiveId string) (urlAmount int, status string, downloadUrl interface{}, error api_error.IApiError)
	CreateArchive() (id string, error api_error.IApiError)
	AddUrlToArchive(archiveId string, url string) api_error.IApiError
	GetArchivePath(archiveId string) (string, api_error.IApiError)
}

type ArchiveService struct {
	storage archiveStorage.IArchiveStorage
	config  *config.Config
}

func NewArchiveService(storage archiveStorage.IArchiveStorage, config *config.Config) IArchiveService {
	return &ArchiveService{
		storage: storage,
		config:  config,
	}
}

func (a *ArchiveService) GetArchiveStatus(archiveId string) (urlAmount int, status string, downloadUrl interface{}, error api_error.IApiError) {
	if !a.storage.HasArchive(archiveId) {
		return 0, "", nil, &api_error.NotFound{Message: ArchiveNotFound}
	}

	amount := a.storage.GetAmountOfUrlsInArchive(archiveId)
	if amount == a.config.ObjectsInArchiveLimit {
		url := a.config.Scheme + "://" + a.config.Hostname + ":" + strconv.Itoa(a.config.Port) + "/api/archive/" + archiveId + "/download"
		return amount, a.storage.GetStatus(archiveId), url, nil
	}

	return amount, a.storage.GetStatus(archiveId), nil, nil
}

func (a *ArchiveService) CreateArchive() (id string, error api_error.IApiError) {

	if a.storage.GetProcessingArchivesAmount() >= a.config.ArchivesLimit {
		return "", &api_error.BadRequest{Message: TooMuchArchives}
	}

	archiveID, err := a.storage.AddArchive()
	if err != nil {
		return "", &api_error.InternalError{}
	}

	return archiveID, nil
}

func (a *ArchiveService) AddUrlToArchive(archiveId string, url string) api_error.IApiError {

	if !a.storage.HasArchive(archiveId) {
		return &api_error.NotFound{Message: ArchiveNotFound}
	}

	if a.storage.GetAmountOfUrlsInArchive(archiveId) >= a.config.ObjectsInArchiveLimit {
		return &api_error.BadRequest{Message: TooMuchObjects}
	}

	// долго выполняется потому что ходит по url для проверки по mime
	if !helpers.IsFileRightFormat(url, a.config.AvailableMimeTypes) {
		return &api_error.UnsupportedMediaType{Message: UnsupportedMediaType}
	}

	a.storage.AddUrl(archiveId, url)

	// если достигли лимита, собираем архив
	if a.storage.GetAmountOfUrlsInArchive(archiveId) == a.config.ObjectsInArchiveLimit {
		go func() {
			err := fs.CreateArchive(archiveId, a.storage.GetUrlList(archiveId), a.config.ArchiveStorageDirPath)
			if err != nil {
				fmt.Print("Не удалось создать архив")
				a.storage.SetStatus(archiveId, archiveStorage.ArchiveStatusError)
			}
			a.storage.SetStatus(archiveId, archiveStorage.ArchiveStatusDone)
		}()
	}
	return nil
}

func (a *ArchiveService) GetArchivePath(archiveId string) (string, api_error.IApiError) {
	filePath := a.config.ArchiveStorageDirPath + archiveId + ".zip"

	if !fs.FileExists(filePath) {
		return "", &api_error.NotFound{Message: ArchiveNotFound}
	}

	return filePath, nil
}
