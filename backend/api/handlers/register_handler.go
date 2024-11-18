package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
func ApiRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Println("Invalid request payload. Got: ", req)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := services.CreateUser(database.DBConnection, &user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		fmt.Println("Error creating user: ", err)
		return
	}

	// if errorMsg != "" {
	// 	data := types.ResonseObject{
	// 		Error: errorMsg,
	// 		Form: map[string]string{
	// 			"Username": username,
	// 			"Email":    email,
	// 		},
	// 	}
	// 	tmpl.ExecuteTemplate(w, "register.html", data)
	// 	return
	// }
}
