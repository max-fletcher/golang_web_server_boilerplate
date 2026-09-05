package fileupload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func LocalFileUploader(r *http.Request, storeFiles []string, destination string, baseUrl string, optional bool) (map[string]string, error) {
	storedFiles := map[string]string{}
	for _, storeFile := range storeFiles {
		file, fileHeader, err := r.FormFile(storeFile)
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) { // We would need an else condition if we made file required
				if optional { // if file is optional, just continue
					continue
				}
				RollbackFiles(storedFiles)
				return nil, ErrFileStorageError{Err: fmt.Errorf("Required file missing: %v", storeFile)}
			}

			return nil, ErrFileStorageError{Err: fmt.Errorf("Failed to upload file")}
		}

		// Set filename and destination/path
		filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
		pathToFile := filepath.Join("uploads", destination, filename)

		if err := os.MkdirAll(filepath.Dir(pathToFile), 0755); err != nil { // Create folder if it doesn't exist
			file.Close()
			RollbackFiles(storedFiles)
			return nil, ErrFileStorageError{
				Err: err,
			}
		}

		// Create or truncate the destination file for writing
		dst, err := os.Create(pathToFile)
		if err != nil {
			file.Close()
			return nil, ErrFileStorageError{Err: err}
		}

		// Store file in destination. It does so by streaming data
		// from the source reader to the destination writer
		_, copyErr := io.Copy(dst, file)
		closeErr := dst.Close()
		fileCloseErr := file.Close()

		if copyErr != nil {
			os.Remove(pathToFile)
			RollbackFiles(storedFiles)

			return nil, ErrFileStorageError{
				Err: copyErr,
			}
		}

		if closeErr != nil {
			os.Remove(pathToFile)
			RollbackFiles(storedFiles)

			return nil, ErrFileStorageError{
				Err: closeErr,
			}
		}

		if fileCloseErr != nil {
			os.Remove(pathToFile)
			RollbackFiles(storedFiles)

			return nil, ErrFileStorageError{
				Err: fileCloseErr,
			}
		}

		fileURL := strings.TrimRight(baseUrl, "/") + "/" + filepath.ToSlash(pathToFile)
		storedFiles[storeFile] = fileURL
	}

	return storedFiles, nil
}

func GetFileHeaders(r *http.Request, storeFiles []string, optional bool) (map[string]*multipart.FileHeader, error) {
	// storedFileHeaders := []*multipart.FileHeader{}
	storedFileHeaders := map[string]*multipart.FileHeader{}

	for _, storeFile := range storeFiles {
		_, fileHeader, err := r.FormFile(storeFile)
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) { // We would need an else condition if we made file required
				if optional { // if file is optional, just continue
					continue
				}
				return nil, ErrFileStorageError{Err: fmt.Errorf("Required file missing: %v", storeFile)}
			}

			return nil, ErrFileStorageError{Err: fmt.Errorf("Failed to upload file")}
		}
		storedFileHeaders[storeFile] = fileHeader
	}

	return storedFileHeaders, nil
}

func RollbackFiles(files map[string]string) {

	for _, file := range files {
		_ = os.Remove(file)
	}
}
