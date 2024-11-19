package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"whoKnows/api/services"
	"whoKnows/database"
	"whoKnows/helperTypes"
	"whoKnows/models"

	"gorm.io/gorm"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	lang := r.URL.Query().Get("language")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}
	if lang == "" {
		lang = "en"
	}
	searchLog := logSearch(query)
	fmt.Println("Query:", query)

	var pages []models.PageData
	results := database.Connection.Where("language = ? AND title ILIKE ?", lang, "%"+query+"%").Order("title")
	if err := results.Find(&pages).Error; err != nil {
		http.Error(w, "Error getting search results", http.StatusInternalServerError)
		return
	}

	body := helperTypes.SearchResponse{
		Data: make([]map[string]interface{}, len(pages)),
	}
	for i, page := range pages {
		body.Data[i] = map[string]interface{}{
			"id":       page.ID,
			"content":  page.Content,
			"language": page.Language,
			"title":    page.Title,
			"url":      page.Url,
			"snippet":  getSnippet(page.Content, query),
		}
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Search results",
		"data":    body,
		"results": len(pages),
		"hits":    searchLog.Count + 1,
	}

	services.ResponseSuccess(w, response, http.StatusOK)
}

func logSearch(query string) models.SearchLog {
	db := database.Connection
	var searchLog models.SearchLog
	if err := db.Where("query = ?", query).First(&searchLog).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			services.CreateSearchLog(db, &models.SearchLog{Query: query})
			fmt.Println("Record does not exist")
			return models.SearchLog{Query: query}
		} else {
			fmt.Println("Error:", err.Error())
			return models.SearchLog{}
		}
	} else {
		searchLog.Count++
		db.Save(&searchLog)
		fmt.Println("Record exists")
	}
	fmt.Println("Search existing, incrementing count")
	return searchLog
}

// Needs improvement but good enough for now
func getSnippet(content, query string) string {
	index := strings.Index(strings.ToLower(content), strings.ToLower(query))
	if index == -1 {
		return ""
	}
	start := index
	end := start + 100
	if end > len(content) {
		end = len(content)
	}
	return content[start:end]
}
