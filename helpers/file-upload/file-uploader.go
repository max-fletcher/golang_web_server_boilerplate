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

func LocalFileUploader(r *http.Request, filenamesToStore []string, destination string, baseUrl string, optional bool) (map[string]string, error) {
	storedFiles := map[string]string{}
	for _, storeFile := range filenamesToStore {
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

// Got rid of "optional" param since we don't need "optional" to throw error anymore. Validation will fail instead.
// func GetFileHeaders(r *http.Request, filenamesToStore []string, optional bool) (map[string]*multipart.FileHeader, error) {
func GetFileHeaders(r *http.Request, filenamesToStore []string) (map[string]*multipart.FileHeader, error) {
	storedFileHeaders := map[string]*multipart.FileHeader{}

	// for _, storeFile := range filenamesToStore {
	// 	_, fileHeader, err := r.FormFile(storeFile)
	// 	if err != nil {
	// 		if errors.Is(err, http.ErrMissingFile) {
	// 			continue
	// 			// replaced lines below with "continue" above because I don't want this to return error here. Instead, I want validation to fail
	// 			// if file is mandatory but doesn't exist
	// 			// if optional { // if file is optional, just continue
	// 			// 	continue
	// 			// }
	// 			// return nil, ErrFileStorageError{Err: fmt.Errorf("Required file missing: %v", storeFile)}
	// 		}

	// 		return nil, ErrFileStorageError{Err: fmt.Errorf("Failed to upload file")}
	// 	}
	// 	storedFileHeaders[storeFile] = fileHeader
	// }

	// Replaced the block above with the block below as we don't need to get entire file toget file header. Works if
	// r.ParseMultipartForm(...) was previously used
	for _, fileName := range filenamesToStore {
		headers := r.MultipartForm.File[fileName]

		if len(headers) == 0 {
			continue
		}

		storedFileHeaders[fileName] = headers[0]
	}

	return storedFileHeaders, nil
}

func RollbackFiles(files map[string]string) {

	for _, file := range files {
		_ = os.Remove(file)
	}
}
