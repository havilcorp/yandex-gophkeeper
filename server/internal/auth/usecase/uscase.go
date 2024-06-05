package usecase

import (
	"crypto/md5"
	"encoding/hex"

	"ya-gophkeeper-server/internal/auth"
	"ya-gophkeeper-server/internal/auth/entity"
)

type usecace struct {
	repo auth.Repository
}

func New(repo auth.Repository) *usecace {
	return &usecace{
		repo: repo,
	}
}

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

// TODO: check GetUser error
func (uc *usecace) Registration(email string, password string) (*entity.User, error) {
	_, err := uc.repo.GetUser(email)
	if err == nil {
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
