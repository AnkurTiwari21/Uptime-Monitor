package utils

import (
	"fmt"
	"time"

	"github.com/ankur12345678/uptime-monitor/pkg/constants"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(secret string, email string, expiresIn int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"jti":   uuid.New().String(),
			"exp":   time.Now().Add(time.Second * time.Duration(expiresIn)).Unix(),
		})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func UUIDGen(category string) string {
	id, err := gonanoid.New(10)
	if err != nil {
		logger.Error("Error in generating uuid ", err)
	}
	switch category {
	case constants.USER_TYPE:
		return fmt.Sprintf("user_%s", id)
	case constants.WEBISTE_TYPE:
		return fmt.Sprintf("web_%s", id)
	case constants.INCIDENT_EVENT_TYPE:
		return fmt.Sprintf("ie_%s", id)
	}
	return ""
}
