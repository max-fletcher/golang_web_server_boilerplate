package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
)

func DecodeJSON(w http.ResponseWriter, r *http.Request, payload any) bool {
	// decoder := json.NewDecoder(r.Body) // decode request body
	// err := decoder.Decode(&params)     // decode request body(json) and store in params variable, or get an error
	err := json.NewDecoder(r.Body).Decode(payload) // This can be used instead of the above 2 lines if you want 1 single line instead of 2
	if err == nil {
		return true
	}

	var typeError *json.UnmarshalTypeError // temp var for storing detailed info about error
	if errors.As(err, &typeError) {
		errormap := map[string]string{typeError.Field: typeError.Field + " must be on type " + typeError.Type.String()}
		responses.RespondWithDetailedErrors(w, http.StatusBadRequest, "Invalid field type", errormap)
		return false
	}

	responses.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))

	return false
}
