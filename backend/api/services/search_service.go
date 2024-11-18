package services

import (
	"fmt"
	"whoKnows/models"

	"gorm.io/gorm"
)

func CreateSearchLog(db *gorm.DB, searchLog *models.SearchLog) error {
	err := db.Create(searchLog).Error

	if err != nil {
		fmt.Printf("error creating search log. Error: %s. Search: %s", err, searchLog.Query)
	}
	return nil
}

func GetSearchLog(db *gorm.DB, search string) (*models.SearchLog, error) {
	var searchLog models.SearchLog
	err := db.Where("search = ?", search).First(&searchLog).Error

	if err != nil {
		return nil, fmt.Errorf("error getting search log. Error: %s. Search: %s", err, search)
	}
	return &searchLog, nil
}

func CreatePageData(db *gorm.DB, pageData *models.PageData) error {
	err := db.Create(pageData).Error

	if err != nil {
		fmt.Printf("error creating page data. Error: %s. PageData: %s", err, pageData.Url)
	}
	return nil
}

func GetPageDataByUrl(db *gorm.DB, url string) (*models.PageData, error) {
	var pageData models.PageData
	err := db.Where("url = ?", url).First(&pageData).Error

	if err != nil {
		return nil, fmt.Errorf("error getting page data by url. Error: %s. URL: %s", err, url)
	}
	return &pageData, nil
}

func UpdatePageData(db *gorm.DB, pageData *models.PageData) error {
	err := db.Save(pageData).Error

	if err != nil {
		fmt.Printf("error updating page data. Error: %s. PageData: %s", err, pageData.Url)
	}
	return nil
}
