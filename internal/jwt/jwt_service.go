package jwt

import (
	"fmt"
	"time"

	"github.com/adi-kmt/ai-streak-backend-go/internal/messages"
	"github.com/golang-jwt/jwt/v5"
)

var key = []byte("SecretYouShouldHide")

func GenerateToken(username string) (string, *messages.AppError) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	s, err := t.SignedString(key)
	if err != nil {
		return "", messages.InternalServerError("Unsable to sign token")
	}
	return s, nil
}

func ParseTokenAndGetClaims(tokenString string) (string, *messages.AppError) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(token.Raw), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if usernameClaim, ok := claims["username"].(string); ok {
			return usernameClaim, nil
		} else {
			return "", messages.Forbidden("Invalid token")
		}
	}
	return "", messages.Forbidden("Invalid token")
}
