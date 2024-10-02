package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var ENV_MYSQL_USER, _ = os.LookupEnv("ENV_MYSQL_USER")
var ENV_MYSQL_PASSWORD, _ = os.LookupEnv("ENV_MYSQL_PASSWORD")
var ENV_INIT_MODE, _ = os.LookupEnv("ENV_INIT_MODE")
var DATABASE_PATH = ENV_MYSQL_USER + ":" + ENV_MYSQL_PASSWORD + "@(mysql_db:3306)/whoknows"

type session struct {
	userID   int
	username string
	expiry   time.Time
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
// Run the server on port 8000
func main() {
	fmt.Println("Starting server on port 8080")

	initDB(ENV_INIT_MODE == "true")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/search", searchHandler)
	mux.HandleFunc("/api/login", loginHandler)
	mux.HandleFunc("/api/register", apiRegister)
	mux.HandleFunc("/api/weather", apiWeather)

	// Apply CORS middleware
	handler := corsMiddleware(mux)

	http.ListenAndServe(":8080", handler)
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
	// enableCors(&w)
	db, err := sql.Open("mysql", DATABASE_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	language := r.FormValue("language")
	if language == "" {
		language = "en"
	}

	rows, err := queryDB(db, "SELECT * FROM pages WHERE language = ? AND content LIKE ?", language, "%"+query+"%")
	if err != nil {
		log.Fatal(err)
	}

	response := map[string]interface{}{
		"data": rows,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// enableCors(&w)

	log.Printf("Received %s request for /api/login", r.Method)

	db, err := sql.Open("mysql", DATABASE_PATH)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	var user struct {
		ID       int
		Username string
		Password string
	}

	err = db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	if !verifyPassword(user.Password, password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New().String()
	sessionStore[sessionID] = session{userID: user.ID, username: username, expiry: time.Now().Add(24 * time.Hour)}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
	})

	println(sessionID)

	response := map[string]interface{}{
		"message":    "Logged in successfully",
		"statusCode": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func verifyPassword(storedHash, password string) bool {
	hash := md5.Sum([]byte(password))
	return storedHash == hex.EncodeToString(hash[:])
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

func apiRegister(w http.ResponseWriter, r *http.Request) {
    // enableCors(&w)
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

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
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
        return
    }
}

func apiWeather(w http.ResponseWriter, r *http.Request) {
	// Handle the actual weather API logic here
	// For now, return a dummy response
	fmt.Fprintf(w, "Weather API is not implemented yet.")
}