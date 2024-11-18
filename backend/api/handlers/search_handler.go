package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"whoKnows/models"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl = template.Must(template.ParseFiles(
	"../app/frontend/root.html",
	"../app/frontend/search.html",
	"../app/frontend/register.html",
	"../app/frontend/login.html",
	"../app/frontend/about.html",
))

func ApiSearchHandler(w http.ResponseWriter, r *http.Request) {
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

	var searchResults []models.PageData
	if query != "" {
		rows, err := QueryDB(db, "SELECT * FROM pages WHERE language = ? AND content LIKE ?", language, "%"+query+"%")
		if err != nil {
			log.Fatal(err)
		}

		searchResults = make([]models.PageData, len(rows))
		for i, row := range rows {
			url, ok := row["url"].(string)
			if !ok {
				url = ""
			}
			title, ok := row["title"].(string)
			if !ok {
				title = ""
			}
			content, ok := row["content"].(string)
			if !ok {
				content = ""
			}

			searchResults[i] = models.PageData{
				Url:     url,
				Title:   title,
				Content: content,
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"search_results": searchResults,
	})
}

func QueryDB(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
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
