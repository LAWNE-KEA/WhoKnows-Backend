package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
)


// Run the server on port 8080
func main() {
	fmt.Println("Starting server on port 8000")

	initDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})


	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := queryDB(db, "SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "Test data: %v", rows)

	})


	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
			return
		}

		rows, err := queryDB(db, "SELECT * FROM pages WHERE content LIKE %?%")
		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Fprintf(w, "Search results: %v", rows)
	})


	http.ListenAndServe(":8000", nil)
}

var DATABASE_PATH = "./tmp/whoknows.db"

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

func queryDB(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
    rows, err := db.Query(query, args...)
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