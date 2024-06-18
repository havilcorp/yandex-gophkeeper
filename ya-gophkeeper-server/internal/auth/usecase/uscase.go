// Package usecase пакет бизнес логики авторизации
package usecase

import (
	"crypto/md5"
	"encoding/hex"

	"yandex-gophkeeper-server/internal/auth"
	"yandex-gophkeeper-server/internal/auth/entity"
)

type usecace struct {
	repo auth.Repository
}

// New получить экземпляр структуры
func New(repo auth.Repository) *usecace {
	return &usecace{
		repo: repo,
	}
}

// Login авторизация пользователя
func (uc *usecace) Login(email string, password string) (*entity.User, error) {
	user, err := uc.repo.GetUser(email)
	if err != nil {
		return nil, entity.ErrUserNotFound
	}
	hash := md5.Sum([]byte(password))
	hashPassword := hex.EncodeToString(hash[:])
	if user.Password != hashPassword {
		return nil, entity.ErrWrongPassword
	}
	return user, nil
}

// Registration регистрация пользователя
func (uc *usecace) Registration(email string, password string) (*entity.User, error) {
	if _, err := uc.repo.GetUser(email); err == nil {
		return nil, entity.ErrUserExists
	}
	hash := md5.Sum([]byte(password))
	hashPassword := hex.EncodeToString(hash[:])
	user, err := uc.repo.CreateUser(email, hashPassword)
	if err != nil {
		return nil, err
	}
	return user, nil
}
