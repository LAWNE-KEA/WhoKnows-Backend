package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var tmpl = template.Must(template.ParseFiles(
	"../frontend/layout.html",
	"../frontend/search.html",
	"../frontend/register.html",
	"../frontend/login.html",
	"../frontend/about.html",
))

var ENV_MYSQL_USER, _ = os.LookupEnv("ENV_MYSQL_USER")
var ENV_MYSQL_PASSWORD, _ = os.LookupEnv("ENV_MYSQL_PASSWORD")
var ENV_INIT_MODE, _ = os.LookupEnv("ENV_INIT_MODE")
var DATABASE_PATH = ENV_MYSQL_USER + ":" + ENV_MYSQL_PASSWORD + "@(mysql_db:3306)/whoknows"

type PageData struct {
	User          *User
	Flashes       []string
	Query         string
	SearchResults []SearchResult
	Error         string
	Form          map[string]string
}

type User struct {
	Username string
}

type SearchResult struct {
	URL         string
	Title       string
	Description string
}

type session struct {
	userID   int
	username string
	expiry   time.Time
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

var sessionStore = map[string]session{}

// Middleware to handle CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Run the server on port 8080
func main() {
	fmt.Println("Starting server on port 8080")

	fs := http.FileServer(http.Dir("./static"))
       http.Handle("/static/", http.StripPrefix("/static/", fs))

	initDB(ENV_INIT_MODE == "true")

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/weather.html")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/root.html")
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/search.html")
	})
	mux.HandleFunc("/", searchHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/api/search", apiSearchHandler)
	mux.HandleFunc("/api/login", apiLoginHandler)
	mux.HandleFunc("/api/register", apiRegisterHandler)

	fs := http.FileServer(http.Dir("../frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/about.html")
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/register.html")
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/login.html")
	})

	// Apply CORS middleware
	handler := corsMiddleware(mux)
 
	http.ListenAndServe(":8080", handler)
}

func init() {
	mime.AddExtensionType(".css", "text/css")
}

func initDB(initMode bool) {
	// Open the SQL file
	if initMode {
		sqlFile, err := os.ReadFile("./schema.sql")
		if err != nil {
			log.Fatal(err)
		}

		// Convert SQL bytes to string
		sqlCommands := string(sqlFile)

		// Parse the SQL commands
		commands := parseSQLCommands(sqlCommands)

		// Open the database connection
		db, err := sql.Open("mysql", DATABASE_PATH)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Execute each command separately
		for _, command := range commands {
			// Trim whitespace and skip empty commands
			command = strings.TrimSpace(command)
			if command == "" {
				continue
			}

			_, err := db.Exec(command)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("SQL commands executed successfully")
	}
}

func parseSQLCommands(sqlCommands string) []string {
	var commands []string
	var currentCommand strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for _, char := range sqlCommands {
		switch char {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote {
				commands = append(commands, currentCommand.String())
				currentCommand.Reset()
				continue
			}
		}
		currentCommand.WriteRune(char)
	}

	// Add the last command if any
	if currentCommand.Len() > 0 {
		commands = append(commands, currentCommand.String())
	}

	return commands
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", DATABASE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := r.FormValue("q")
	language := r.FormValue("language")
	if language == "" {
		language = "en"
	}

	var searchResults []SearchResult
	if query != "" {
		rows, err := queryDB(db, "SELECT * FROM pages WHERE language = ? AND content LIKE ?", language, "%"+query+"%")
		if err != nil {
			log.Fatal(err)
		}

		searchResults = make([]SearchResult, len(rows))
		for i, row := range rows {
			searchResults[i] = SearchResult{
				URL:         row["url"].(string),
				Title:       row["title"].(string),
				Description: row["description"].(string),
			}
		}
	}

	data := PageData{
		Query:         query,
		SearchResults: searchResults,
	}

	tmpl.ExecuteTemplate(w, "search.html", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		User:    &User{Username: "JohnDoe"}, // Example user, replace with actual user data
		Flashes: []string{"Welcome to the About page!"},
	}
	tmpl.ExecuteTemplate(w, "about.html", data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		User:    &User{Username: "JohnDoe"}, // Example user, replace with actual user data
		Flashes: []string{"Please log in."},
	}
	tmpl.ExecuteTemplate(w, "login.html", data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		User:    &User{Username: "JohnDoe"}, // Example user, replace with actual user data
		Flashes: []string{"Please register."},
	}
	tmpl.ExecuteTemplate(w, "register.html", data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No active session", http.StatusBadRequest)
		return
	}

	delete(sessionStore, sessionID.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func apiSearchHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", DATABASE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := r.FormValue("q")
	language := r.FormValue("language")
	if language == "" {
		language = "en"
	}

	var searchResults []SearchResult
	if query != "" {
		rows, err := queryDB(db, "SELECT * FROM pages WHERE language = ? AND content LIKE ?", language, "%"+query+"%")
		if err != nil {
			log.Fatal(err)
		}

		searchResults = make([]SearchResult, len(rows))
		for i, row := range rows {
			searchResults[i] = SearchResult{
				URL:         row["url"].(string),
				Title:       row["title"].(string),
				Description: row["description"].(string),
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"search_results": searchResults,
	})
}

func apiLoginHandler(w http.ResponseWriter, r *http.Request) {
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

	if !verifyPassword(storedHash, password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New().String()
	expirationTime := time.Now().Add(24 * time.Hour)
	sessionStore[sessionID] = session{userID: 1, username: username, expiry: expirationTime}

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

func apiRegisterHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", DATABASE_PATH)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	var errorMsg string
	if username == "" {
		errorMsg = "You have to enter a username"
	} else if email == "" || !strings.Contains(email, "@") {
		errorMsg = "You have to enter a valid email address"
	} else if password == "" {
		errorMsg = "You have to enter a password"
	} else if password != password2 {
		errorMsg = "The two passwords do not match"
	} else {
		// Check if username already exists
		_, err := getUserID(db, username)
		if err == nil {
			errorMsg = "The username is already taken"
		} else if err != sql.ErrNoRows {
			// An unexpected error occurred
			log.Printf("Error checking username: %v", err)
			http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
			return
		} else {
			// Username is available, proceed with registration
			hashedPassword := hashPassword(password)
			_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
			if err != nil {
				log.Printf("Error inserting new user: %v", err)
				http.Error(w, "Failed to register user", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Registered successfully"})
			return
		}
	}

	if errorMsg != "" {
		data := PageData{
			Error: errorMsg,
			Form: map[string]string{
				"Username": username,
				"Email":    email,
			},
		}
		tmpl.ExecuteTemplate(w, "register.html", data)
		return
	}
}

func queryDB(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(query, args...)
	fmt.Println("Query:", query, "Args:", args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				row[colName] = string(b)
			} else {
				row[colName] = val
			}
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func getUserID(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func verifyPassword(storedHash, password string) bool {
	hash := md5.Sum([]byte(password))
	return storedHash == hex.EncodeToString(hash[:])
}
