package security

import (
	"fmt"
	"time"
	"whoKnows/api/configs"
	"whoKnows/models"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func CreateJWT(userId uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"iss":      "whoKnows",
		"sub":      userId,
		"aud":      "whoKnows",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
		"username": username,
		"role":     "user",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := configs.EnvConfig.JWT.Secret

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	secret := configs.EnvConfig.JWT.Secret

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// not finished so doesnt work correctly yet
func ExpireJWT(database *gorm.DB, tokenString string) error {
	token := &models.Token{Token: tokenString}
	token.ExpiresAt = time.Now()
	err := database.Create(token).Error
	if err != nil {
		fmt.Printf("error invalidating token. Error: %s. Token: %s", err, tokenString)
	}

	fmt.Printf("token invalidated: %s", tokenString)
	return nil
}

// func VerifyExpiredJWT(database *gorm.DB, tokenString string) (bool, error) {
// 	var token models.Token
// 	err := database.Where("token = ?", tokenString).First(&token).Error
// 	if err != nil {
// 		return false, fmt.Errorf("Error finding token. Error: %s. Token: %s", err, tokenString)
// 	}

// 	if token.ExpiresAt.Before(time.Now()) {
// 		return true, nil
// 	}

// 	return false, nil
// }
