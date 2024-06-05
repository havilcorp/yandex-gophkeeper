package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var jwtkey = []byte("very-secret-key") // TODO: JWT secret

func GenerateJWT(userId string) (string, error) {
	payload := jwt.MapClaims{
		"sub": userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtkey)
	if err != nil {
		return "", fmt.Errorf("token signed string: %w", err)
	}
	return t, nil
}

func VerifyJWT(token string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(jwtkey), nil
	}
	payload := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &payload, keyFunc)
	if err != nil {
		return "", fmt.Errorf("parse with claim: %w", err)
	}
	userId, err := payload.GetSubject()
	if err != nil {
		return "", fmt.Errorf("get subject: %w", err)
	}
	return userId, nil
}
