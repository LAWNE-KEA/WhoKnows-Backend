package services

import (
	"fmt"
	"whoKnows/models"

	"gorm.io/gorm"
)

func CreateSearchLog(db *gorm.DB, searchLog *models.SearchLog) error {
	err := db.Create(searchLog).Error

	if err != nil {
		fmt.Printf("error creating search log. Error: %s. SearchLog: %s", err, searchLog.Search)
	}
	return nil
}
