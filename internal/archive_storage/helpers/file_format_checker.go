package helpers

import (
	"mime"
	"net/http"
	"strings"
)

func IsFileRightFormat(fileUrl string, mimeTypes []string) bool {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		mt, _, err := mime.ParseMediaType(contentType)
		if err == nil {
			contentType = strings.ToLower(mt)
			for _, allowed := range mimeTypes {
				if contentType == strings.ToLower(allowed) {
					return true
				}
			}
			return false
		}
	}

	buffer := make([]byte, 512)
	n, err := resp.Body.Read(buffer)
	if err != nil && n == 0 {
		return false
	}

	detected := strings.ToLower(http.DetectContentType(buffer[:n]))
	for _, allowed := range mimeTypes {
		if detected == strings.ToLower(allowed) {
			return true
		}
	}

	return false
}
