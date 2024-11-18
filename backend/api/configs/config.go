package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var ENV_MYSQL_USER, _ = os.LookupEnv("ENV_MYSQL_USER")
var ENV_MYSQL_PASSWORD, _ = os.LookupEnv("ENV_MYSQL_PASSWORD")
var ENV_INIT_MODE, _ = os.LookupEnv("ENV_INIT_MODE")
var DATABASE_PATH = ENV_MYSQL_USER + ":" + ENV_MYSQL_PASSWORD + "@(mysql_db:3306)/whoknows"

var EnvConfig Config

func LoadEnv() error {
	path := os.Getenv("ENV_PATH")

	if path != "" {
		if err := godotenv.Load(path); err != nil {
			return fmt.Errorf("error loading environment variables: %s, path: %s", err, path)
		}
	}
	fmt.Println("Loaded environment variables")
	return nil
}

func SetEnv(key, value string) error {
	if err := os.Setenv(key, value); err != nil {
		return fmt.Errorf("error setting environment variable %s: %s", key, err)
	}
	return nil
}

func GetEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s not found", key)
}

type Config struct {
	Database struct {
		Host     string
		User     string
		Password string
		Name     string
		Port     int
		SSLMode  string
		Migrate  bool
	}
}
