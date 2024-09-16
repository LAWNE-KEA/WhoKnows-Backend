package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_"github.com/mattn/go-sqlite3"
    "io/ioutil"
    "log"
    "encoding/json"
    "crypto/md5"
    "encoding/hex"
    "strings"
)

var DATABASE_PATH = "./tmp/whoknows.db"


// Run the server on port 8000
func main() {
	fmt.Println("Starting server on port 8000")

	initDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

    http.HandleFunc("/api/search", searchHandler)
    
    http.HandleFunc("/api/login", loginHandler)

    http.HandleFunc("/api/register", apiRegister)


	http.ListenAndServe(":8000", nil)
}


func initDB() {
    db, err := sql.Open("sqlite3", DATABASE_PATH)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    schema, err := ioutil.ReadFile("./schema.sql")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(string(schema))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Initialized the database:", DATABASE_PATH)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
  
    db, err := sql.Open("sqlite3", DATABASE_PATH)
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

    fmt.Fprintf(w, "Search results: %v", rows)
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

	db, err := sql.Open("sqlite3", DATABASE_PATH)
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
    
	// Here you would typically create a session or generate a token
	// For this example, we'll just return a success message
    // Generate a simple token for demonstration purposes
    token := fmt.Sprintf("token-%d", user.ID)
    w.Header().Set("Authorization", token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged in successfully", "user_id": string(user.ID)})
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
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    db, err := sql.Open("sqlite3", DATABASE_PATH)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var error string
    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")
    password2 := r.FormValue("password2")

    if username == "" {
        error = "You have to enter a username"
    } else if email == "" || !strings.Contains(email, "@") {
        error = "You have to enter a valid email address"
    } else if password == "" {
        error = "You have to enter a password"
    } else if password != password2 {
        error = "The two passwords do not match"
    } else if _, err := getUserID(db, username); err == nil {
        error = "The username is already taken"
    } else {
        hashedPassword := hashPassword(password)
        _, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Fprintf(w, "Registered successfully")
        return
    }

    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"error": error})
}