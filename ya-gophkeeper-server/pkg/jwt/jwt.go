package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(key []byte, user string) (string, error) {
	payload := jwt.MapClaims{
		"sub": user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("token signed string: %w", err)
	}
	return t, nil
}

func VerifyJWT(key []byte, token string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(key), nil
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
