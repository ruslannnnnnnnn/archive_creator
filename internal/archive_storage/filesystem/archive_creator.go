package filesystem

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateArchive(name string, urlList []string, path string) error {

	outputPath := filepath.Join(path, name+".zip")

	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка создания архива: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, url := range urlList {
		func() {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("ошибка загрузки %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("некорректный статус %s: %d\n", url, resp.StatusCode)
				return
			}

			segments := strings.Split(url, "/")
			filename := segments[len(segments)-1]
			if filename == "" {
				filename = fmt.Sprintf("file_%d", time.Now().UnixNano())
			}

			zipFileWriter, err := zipWriter.Create(filename)
			if err != nil {
				fmt.Printf("ошибка создания zip-файла %s: %v\n", filename, err)
				return
			}

			_, err = io.Copy(zipFileWriter, resp.Body)
			if err != nil {
				fmt.Printf("ошибка записи файла %s в архив: %v\n", filename, err)
				return
			}
		}()
	}

	fmt.Printf("Архив успешно создан: %s\n", outputPath)
	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
