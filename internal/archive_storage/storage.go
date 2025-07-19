package archive_storage

import (
	"github.com/google/uuid"
	"sync"
)

const ArchiveStatusInProgress = "In Progress"
const ArchiveStatusDone = "Done"

// Archive задача на создание архива
type Archive struct {
	FileUrls []string
	Status   string
}

// Storage хранилище задач на создание архива
type Storage struct {
	archives map[string]*Archive
	mu       sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		archives: make(map[string]*Archive),
	}
}

func (s *Storage) GetProcessingArchivesAmount() int {
	s.mu.Lock()
	result := 0
	for _, v := range s.archives {
		if v.Status == ArchiveStatusInProgress {
			result++
		}
	}
	s.mu.Unlock()
	return result
}

func (s *Storage) GetAmountOfUrlsInArchive(archiveId string) int {
	s.mu.Lock()
	result := len(s.archives[archiveId].FileUrls)
	s.mu.Unlock()
	return result
}

func (s *Storage) HasArchive(id string) bool {
	s.mu.Lock()
	_, result := s.archives[id]
	s.mu.Unlock()
	return result
}

func (s *Storage) GetUrlList(id string) []string {
	s.mu.Lock()
	result := s.archives[id].FileUrls
	s.mu.Unlock()
	return result
}

func (s *Storage) AddArchive() (id string, err error) {
	s.mu.Lock()

	archiveId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	s.archives[archiveId.String()] = &Archive{
		FileUrls: make([]string, 0),
		Status:   ArchiveStatusInProgress,
	}
	s.mu.Unlock()

	return archiveId.String(), nil
}

func (s *Storage) AddUrl(archiveId string, url string) {
	s.mu.Lock()
	s.archives[archiveId].FileUrls = append(s.archives[archiveId].FileUrls, url)
	s.mu.Unlock()
}

func (s *Storage) GetStatus(id string) string {
	s.mu.Lock()
	result := s.archives[id].Status
	s.mu.Unlock()
	return result
}

func (s *Storage) SetStatus(id string, status string) {
	s.mu.Lock()
	s.archives[id].Status = status
	s.mu.Unlock()
}
