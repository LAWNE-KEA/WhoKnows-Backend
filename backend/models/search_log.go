package models

import "time"

type SearchLog struct {
	ID        uint   `gorm:"primaryKey"`
	Search    string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
}
