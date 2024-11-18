package services

import (
	"errors"
	"fmt"
	"time"
	"whoKnows/models"
	"whoKnows/security"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	err := db.Create(user).Error

	if err != nil {
		fmt.Printf("error creating user. Error: %s. User: %s", err, user.Username)
	}
	return nil
}

func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("error getting user by username. Error: %s. Username: %s", err, username)
	}
	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("error getting user by email. Error: %s. Email: %s", err, email)
	}
	return &user, nil
}

func GetUserByID(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	err := db.Where("id = ?", id).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("error getting user by id. Error: %s. ID: %d", err, id)
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *models.User) error {
	err := db.Save(user).Error

	if err != nil {
		fmt.Printf("error updating user. Error: %s. User: %s", err, user.Username)
	}
	return nil
}

func DeleteUser(db *gorm.DB, user *models.User) error {
	err := db.Delete(user).Error

	if err != nil {
		fmt.Printf("error deleting user. Error: %s. User: %s", err, user.Username)
	}
	return nil
}

func UpdateLastActive(db *gorm.DB, user *models.User) error {
	err := db.Model(user).Update("last_active", time.Now()).Error

	if err != nil {
		fmt.Printf("error updating last active. Error: %s. User: %s", err, user.Username)
	}
	return nil
}

func CheckPassword(db *gorm.DB, username, password string) (*models.User, bool, error) {
	user, err := GetUserByUsername(db, username)
	if err != nil {
		return nil, false, errors.New("invalid username or password")
	}
	if !security.VerifyPassword(user.Password, password) {
		return nil, false, errors.New("invalid username or password")
	}

	return user, true, nil
}
