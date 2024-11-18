package models

import "time"

type SearchLog struct {
	ID        uint   `gorm:"primaryKey"`
	Query     string `gorm:"type:varchar(255);not null"`
	Count     int    `gorm:"not null;default:0"`
	CreatedAt time.Time
}
