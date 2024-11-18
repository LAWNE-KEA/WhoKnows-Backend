package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"whoKnows/security"
)

var ENV_MYSQL_USER, _ = os.LookupEnv("ENV_MYSQL_USER")
var ENV_MYSQL_PASSWORD, _ = os.LookupEnv("ENV_MYSQL_PASSWORD")
var ENV_INIT_MODE, _ = os.LookupEnv("ENV_INIT_MODE")
var DATABASE_PATH = ENV_MYSQL_USER + ":" + ENV_MYSQL_PASSWORD + "@(mysql_db:3306)/whoknows"

func ApiLoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", DATABASE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedHash string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedHash)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !security.VerifyPassword(storedHash, password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New().String()
	expirationTime := time.Now().Add(24 * time.Hour)
	//sessionStore[sessionID] = session{userID: 1, username: username, expiry: expirationTime}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Logged in successfully",
		"statusCode": http.StatusOK,
		"username":   username,
		"sessionID":  sessionID,
	})
}
