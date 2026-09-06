package validator

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

func ValidatePhoto(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return errors.New("Invalid file")
	}

	// if file == nil {
	// 	// return nil // use this if file is optional
	// 	return errors.New("Photo is required")
	// }

	if file != nil {
		if file.Size > 5*1024*1024 {
			return errors.New("Photo must not exceed 5 MB")
		}

		src, err := file.Open()
		if err != nil {
			return errors.New("Unable to read photo")
		}
		defer src.Close()

		buffer := make([]byte, 512)

		_, err = src.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return errors.New("Unable to read photo")
		}

		contentType := http.DetectContentType(buffer)

		switch contentType {
		case "image/jpeg", "image/png", "image/webp":
			return nil
		default:
			return errors.New("Photo must be JPEG, PNG, or WebP")
		}
	}

	return nil
}
