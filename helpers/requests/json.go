package requests

import (
	"encoding/json"
	"errors"
	"net/http"
)

func DecodeJSON(r *http.Request, payload any) error {
	// decoder := json.NewDecoder(r.Body) // decode request body
	// err := decoder.Decode(&params)     // decode request body(json) and store in params variable, or get an error
	err := json.NewDecoder(r.Body).Decode(payload) // This can be used instead of the above 2 lines if you want 1 single line instead of 2
	if err == nil {
		return nil
	}

	var typeError *json.UnmarshalTypeError // temp var for storing detailed info about error
	if errors.As(err, &typeError) {
		return ErrInvalidFieldType{
			Field: typeError.Field,
			Type:  typeError.Type.String(),
		}
	}

	return ErrInvalidJSON{
		Err: err,
	}
}
