package helperTypes

import (
	"whoKnows/models"
)

type ResponseData struct {
	User          *models.User
	Flashes       []string
	Query         string
	SearchResults []models.PageData
	Error         string
	Form          map[string]string
}

type SearchResponse struct {
	Data []map[string]interface{} `json:"data"`
}
