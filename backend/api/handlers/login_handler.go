package handlers

import (
	"encoding/json"
	"net/http"

	"whoKnows/api/services"
	"whoKnows/database"
	"whoKnows/security"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		services.ResponseError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		services.ResponseError(w, "username and password are required", http.StatusBadRequest)
		return
	}

	user, check, err := services.CheckPassword(database.Connection, loginRequest.Username, loginRequest.Password)
	if err != nil || !check {
		services.ResponseError(w, "invalid username or password", http.StatusInternalServerError)
		return
	}

	token, err := security.CreateJWT(user.ID, user.Username)
	if err != nil {
		services.ResponseError(w, "error creating token", http.StatusInternalServerError)
		return
	}

	services.ResponseSuccess(w,
		map[string]interface{}{
			"status":  "success",
			"message": "Login successful",
			"token":   token,
		},
		http.StatusOK,
	)

}
