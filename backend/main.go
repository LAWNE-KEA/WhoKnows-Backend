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
    "github.com/google/uuid"
    "time"
    "path/filepath"
    "os"
)

var DATABASE_PATH = "./tmp/whoknows.db"


var sessionStore = map[string]session{}

type session struct {
    userName string
    expiry time.Time
}

// Run the server on port 8000
func main() {
    // Get the current working directory
    cwd, err := os.Getwd()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print the current working directory
    fmt.Println("Current working directory:", cwd)


	fmt.Println("Starting server on port 8080")

	initDB()

    

	http.HandleFunc("/", serveIndexPage)

    http.HandleFunc("/api/search", searchHandler)
    
    http.HandleFunc("/api/login", loginHandler)

    http.HandleFunc("/api/register", apiRegister)

    // Serve static files
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
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
    enableCors(&w)
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

    json.NewEncoder(w).Encode(rows)
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
    enableCors(&w)

    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }


    if (*r).Method == "OPTIONS" {
        return
    }
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
    
    sessionID := uuid.New().String()
    sessionStore[sessionID] = session{userName: username, expiry: time.Now().Add(24 * time.Hour)}
    http.SetCookie(w, &http.Cookie{
        Name: "session_id",
        Value: sessionID,
        Expires: time.Now().Add(24 * time.Hour),
    })
    json.NewEncoder(w).Encode(map[string]string{"statusCode": "200", "message": "Logged in successfully"})
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
    enableCors(&w)
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

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    // Get the current working directory
    cwd, err := os.Getwd()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Print the current working directory for debugging
    fmt.Println("Current working directory:", cwd)

    // Construct the path to the HTML file
    filePath := filepath.Join(cwd, "..", "frontend", "whoknows.html")
    
    // Print the file path for debugging
    fmt.Println("Constructed file path:", filePath)

    // Check if the file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        // Print an error message for debugging
        fmt.Println("File not found:", filePath)
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Serve the file
    http.ServeFile(w, r, filePath)
}