package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var EnvConfig Config

func LoadEnv() error {
	path := os.Getenv("ENV_FILE_PATH")
	// path := "../.env"

	if path != "" {
		if err := godotenv.Load(path); err != nil {
			return fmt.Errorf("error loading environment variables: %s, path: %s", err, path)
		}
	}

	populateEnvConfig()

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
		fmt.Printf("found environment variable %s\n", key)
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s not found", key)
}

func GetEnvInt(key string) (int, error) {
	value, err := GetEnv(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func GetEnvBool(key string) (bool, error) {
	value, err := GetEnv(key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(value)
}

func populateEnvConfig() {
	EnvConfig.Database.Host, _ = GetEnv("ENV_DATABASE_HOST")
	EnvConfig.Database.User, _ = GetEnv("ENV_DATABASE_USER")
	EnvConfig.Database.Password, _ = GetEnv("ENV_DATABASE_PASSWORD")
	EnvConfig.Database.Name, _ = GetEnv("ENV_DATABASE_NAME")
	EnvConfig.Database.Port, _ = GetEnvInt("ENV_DATABASE_PORT")
	EnvConfig.Database.SSLMode, _ = GetEnv("ENV_DATABASE_SSL_MODE")
	EnvConfig.Database.Migrate, _ = GetEnvBool("ENV_DATABASE_MIGRATE")
	EnvConfig.Database.SeedFile, _ = GetEnv("ENV_DATABASE_SEED")
	EnvConfig.JWT.Secret, _ = GetEnv("ENV_JWT_SECRET")
	EnvConfig.JWT.Expiry, _ = GetEnvInt("ENV_JWT_EXPIRY")
	EnvConfig.Weather.APIKey, _ = GetEnv("ENV_WEATHER_API_KEY")
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
		SeedFile string
		Log      LogConfig
	}
	JWT struct {
		Secret string
		Expiry int
	}
	Weather struct {
		APIKey string
	}
	Log struct {
		Level string

		Format string
	}
}

var AppConfig Config

type LogConfig struct {
	Level  string
	Format string
}
