package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

// This function will take a payload(struct) and malshal it into a JSON string that will be sent as bytes of data
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error payload inside RespondWithJSON: %v", payload)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithSuccess(w http.ResponseWriter, code int, message string, data any) {
	response := SuccessResponse{
		Code:    code,
		Status:  "ok",
		Message: message,
		Data:    data,
	}

	RespondWithJSON(w, code, response)
}

// This function will take a message(string) and malshal it into a structured JSON string that will be sent as bytes of data
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= http.StatusInternalServerError {
		log.Println("Responding with a 5XX error:", msg)
	}

	RespondWithJSON(w, code, Response{
		Code:    code,
		Status:  "error",
		Message: msg,
	})
}

func RespondWithDetailedErrors(w http.ResponseWriter, code int, message string, err any) {
	RespondWithJSON(w, http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Status:  "error",
		Message: message,
		Errors:  err,
	})
}

func BadRequestError(w http.ResponseWriter, msg string) { // 400 error
	RespondWithError(w, http.StatusBadRequest, msg)
}

func UnauthorizedError(w http.ResponseWriter, msg string) { // 401 error
	RespondWithError(w, http.StatusUnauthorized, msg)
}

func ForbiddenError(w http.ResponseWriter, msg string) { // 403 error
	RespondWithError(w, http.StatusForbidden, msg)
}

func NotFoundError(w http.ResponseWriter, msg string) { // 404 error
	RespondWithError(w, http.StatusNotFound, msg)
}

func ConflictError(w http.ResponseWriter, msg string) { // 409 error
	RespondWithError(w, http.StatusConflict, msg)
}

func ValidationError(w http.ResponseWriter, msg string) { // 422 error
	RespondWithError(w, http.StatusUnprocessableEntity, msg)
}

func ValidationErrorWithFields(w http.ResponseWriter, err any) { // 422 error
	RespondWithJSON(w, http.StatusUnprocessableEntity, ErrorResponse{
		Code:    http.StatusUnprocessableEntity,
		Status:  "error",
		Message: "Validation failed",
		Errors:  err,
	})
}

func InternalServerError(w http.ResponseWriter, msg string) { // 500 error
	RespondWithError(w, http.StatusInternalServerError, msg)
}

func InternalServerErrorSWW(w http.ResponseWriter) { // 500 error
	InternalServerError(w, "Something went wrong. Please try again.")
}
