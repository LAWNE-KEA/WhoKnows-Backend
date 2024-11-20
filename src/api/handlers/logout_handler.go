package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"whoKnows/database"
	"whoKnows/security"
)

// Checks the users JWT token exists before it is invalidated to log the user out
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	var tokenStart = "Bearer "

	validPrefix := strings.HasPrefix(auth, tokenStart)
	if !validPrefix {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	token := strings.Split(auth, tokenStart)[1]

	err := security.ExpireJWT(database.Connection, token)
	if err != nil {
		http.Error(w, "Error invalidating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Logged out successfully"})
}
