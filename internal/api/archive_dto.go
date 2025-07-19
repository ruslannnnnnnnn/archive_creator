package api

type UpdateArchiveDto struct {
	FileUrl string `json:"file_url" binding:"required"`
}

type CreateArchiveResponse struct {
	ArchiveId string `json:"archive_id"`
}

type UpdateArchiveResponse struct {
	ArchiveId string `json:"archive_id"`
}

type GetArchiveStatusResponse struct {
	Status      string      `json:"status"`
	UrlAmount   int         `json:"url_amount"`
	DownloadUrl interface{} `json:"download_url"`
}
