package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

	if configs.EnvConfig.Database.Seed {
		seedData()
	}

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
		return fmt.Errorf("error migrating database: s", err)
	}

	fmt.Println("Schema migrated successfully")
	return nil
}

func seedData() {
	fmt.Println("Seeding data")
	// Open the JSON file
	fmt.Println("Seed file: ", configs.EnvConfig.Database.SeedFile)
	jsonFile, err := os.Open(configs.EnvConfig.Database.SeedFile)
	if err != nil {
		fmt.Println("error opening JSON file: ", err)
		return
	}
	defer jsonFile.Close()

	// Read the JSON file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error reading JSON file: ", err)
		return
	}

	// Unmarshal the JSON data into a slice of PageData
	var pageDataList []models.PageData
	err = json.Unmarshal(byteValue, &pageDataList)
	if err != nil {
		fmt.Println("error unmarshalling JSON data: ", err)
		return
	}

	// Insert each PageData object into the database
	for _, pageData := range pageDataList {
		if pageData.Language == "" {
			pageData.Language = "en"
		}
		if pageData.Title == "" || pageData.Url == "" || pageData.Content == "" {
			fmt.Println("error seeding data: invalid data")
		} else {
			err = Connection.Create(&pageData).Error
			if err != nil {
				fmt.Println("error seeding data:", err)
			}
		}
	}

	fmt.Println("Data seeded successfully")
}
