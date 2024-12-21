package database

import (
	"testing"
	"whoKnows/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initTestDB() (*gorm.DB, error) {
	// Initialize an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the models for testing
	err = db.AutoMigrate(&models.User{}, &models.PageData{}, &models.Token{}, &models.SearchLog{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestDatabaseMigration(t *testing.T) {
	// Initialize the test DB
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Ensure the migrations are working
	if err := db.AutoMigrate(&models.User{}, &models.PageData{}, &models.Token{}, &models.SearchLog{}); err != nil {
		t.Errorf("migration failed: %v", err)
	}

	// Check if the tables exist
	userTableExists := db.Migrator().HasTable(&models.User{})
	if !userTableExists {
		t.Error("User table does not exist after migration")
	}

	pageDataTableExists := db.Migrator().HasTable(&models.PageData{})
	if !pageDataTableExists {
		t.Error("PageData table does not exist after migration")
	}

	tokenTableExists := db.Migrator().HasTable(&models.Token{})
	if !tokenTableExists {
		t.Error("Token table does not exist after migration")
	}

	searchLogTableExists := db.Migrator().HasTable(&models.SearchLog{})
	if !searchLogTableExists {
		t.Error("SearchLog table does not exist after migration")
	}
}
