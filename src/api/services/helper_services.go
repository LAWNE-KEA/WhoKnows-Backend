package services

import (
	"encoding/json"
	"net/http"
)

func ResponseSuccess(w http.ResponseWriter, body map[string]interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func ResponseError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}
