package helpers

import (
	"net/http"
	"strings"
)

func IsFileRightFormat(fileUrl string, mimeTypes []string) bool {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	buffer := make([]byte, 512)
	n, err := resp.Body.Read(buffer)
	if err != nil && n == 0 {
		return false
	}

	mimeType := strings.ToLower(http.DetectContentType(buffer[:n]))

	for _, mt := range mimeTypes {
		if strings.HasPrefix(mimeType, strings.ToLower(mt)) {
			return true
		}
	}

	return false
}
