package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
)

var DATABASE_PATH = "./tmp/whoknows.db"


// Run the server on port 8080
func main() {
	fmt.Println("Starting server on port 8000")

	initDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

// GET VERSION
	// http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
	// 	db, err := sql.Open("sqlite3", DATABASE_PATH)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer db.Close()
		
	// 	query := r.URL.Query().Get("q")
	// 	if query == "" {
	// 		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
	// 		return
	// 	}

    //     language := r.URL.Query().Get("language")
    //     if language == "" {
    //         language = "en"
    //     }

	// 	rows, err := queryDB(db, "SELECT * FROM pages WHERE language = ? AND content LIKE ?", language, "%"+query+"%")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
		
	// 	fmt.Fprintf(w, "Search results: %v", rows)
	// })

// POST VERSION
    http.HandleFunc("/api/search", searchHandler)
    


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
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    
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


func getUserID(db *sql.DB, username string) (int, error) {
    var id int
    err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func apiLogin(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("sqlite3", DATABASE_PATH)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var error string
    username := r.FormValue("username")
    password := r.FormValue("password")

    user, err := queryDB(db, "SELECT * FROM users WHERE username = ?", username)
    if err != nil {
        log.Fatal(err)
    }

    var storedPassword string
    var userID int
    err = user.Scan(&userID, &username, &storedPassword)
    if err == sql.ErrNoRows {
        error = "Invalid username"
    } else if !verifyPassword(storedPassword, password) {
        error = "Invalid password"
    } else {
        fmt.Fprintf(w, "Logged in successfully")
        return
    }

    w.WriteHeader(http.StatusUnauthorized)
    json.NewEncoder(w).Encode(map[string]string{"error": error})
}

func apiRegister(w http.ResponseWriter, r *http.Request) {
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
        _, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashPassword(password))
        if err != nil {
            log.Fatal(err)
        }
        fmt.Fprintf(w, "Registered successfully")
        return
    }

    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{"error": error})
}