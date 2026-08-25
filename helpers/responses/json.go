package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

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

// This function will take a payload(struct) and malshal it into a JSON string that will be sent as bytes of data
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal json response: %v", payload)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
