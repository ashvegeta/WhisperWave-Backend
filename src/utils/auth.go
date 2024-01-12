package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey []byte = []byte("secretKey")

// Function to hash password and such
func HashText(text string) (string, error) {
	saltedHash, err := bcrypt.GenerateFromPassword([]byte(text), 10)
	return string(saltedHash), err
}

// Function to compare the stored hash with the input text
func CompareTextAndHash(text string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text)) == nil
}

// Function to generate a new JWT token
func GenerateToken(UserName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": UserName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Function to verify a JWT token
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	// check if token is valid
	if !token.Valid {
		return fmt.Errorf("invalid token %s", tokenString)
	}

	return nil
}
