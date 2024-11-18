package database

import (
	"fmt"
	"time"

	"whoKnows/api/configs"
	"whoKnows/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Connection *gorm.DB

// Initializes the db connection
func InitDatabase() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		configs.EnvConfig.Database.Host,
		configs.EnvConfig.Database.User,
		configs.EnvConfig.Database.Password,
		configs.EnvConfig.Database.Name,
		configs.EnvConfig.Database.Port,
		configs.EnvConfig.Database.SSLMode,
	)

	var err error
	Connection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to database: %s", err)
	}

	fmt.Println("Database connected")

	if configs.EnvConfig.Database.Migrate {
		return Migrate()
	}

	return nil
}

// migrates the database schema to the latest version
func Migrate() error {
	m := gormigrate.New(Connection, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: time.Now().Format("20060102"),
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.User{}, &models.PageData{}, &models.Token{}, &models.SearchLog{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.User{}, &models.PageData{}, &models.Token{}, &models.SearchLog{})
			},
		},
	})

	err := m.Migrate()
	if err != nil {
		return fmt.Errorf("error migrating database: %s", err)
	}

	fmt.Println("Schema migrated successfully")
	return nil
}
