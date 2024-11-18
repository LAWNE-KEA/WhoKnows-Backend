package types

import (
	"whoKnows/models"
)

type ResonseObject struct {
	User          *models.User
	Flashes       []string
	Query         string
	SearchResults []models.PageData
	Error         string
	Form          map[string]string
}
