package requests

import "net/http"

func ParseFormdata(r *http.Request) error {
	err := r.ParseMultipartForm(10 << 20) // 10 << 20 means 10 multiplied by 2 to the power of 20 so 10MB in bytes
	if err != nil {
		return ErrParsingFormdataError{
			Err: err,
		}
	}

	return nil
}
