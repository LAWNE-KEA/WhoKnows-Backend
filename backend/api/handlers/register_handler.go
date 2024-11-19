package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"whoKnows/api/services"
	"whoKnows/database"
	"whoKnows/models"
	"whoKnows/security"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Password  string `json:"password" validate:"required,min=8"`
	Password2 string `json:"password2" validate:"required,eqfield=Password"`
	Email     string `json:"email" validate:"required,email"`
}

var validate = validator.New()

// Create a new User in the database
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Invalid request payload. Got: ", req)
		services.ResponseError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		services.ResponseError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		services.ResponseError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := services.CreateUser(database.Connection, &user); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			services.ResponseError(w, "Username or email already in use", http.StatusBadRequest)
			return
		} else {
			services.ResponseError(w, "Error creating user", http.StatusInternalServerError)
			return
		}
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "User created successfully",
		"user":    user.Username,
	}

	services.ResponseSuccess(w, response, http.StatusOK)
}
