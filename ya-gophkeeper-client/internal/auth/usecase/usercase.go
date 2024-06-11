package usecase

import (
	"ya-gophkeeper-client/internal/auth"
	"ya-gophkeeper-client/internal/auth/entity"
)

type UserCase struct {
	del auth.Delivery
}

func New(del auth.Delivery) *UserCase {
	return &UserCase{
		del: del,
	}
}

func (uc *UserCase) Login(dto *entity.LoginDto) (string, error) {
	token, err := uc.del.Login(dto)
	return token, err
}

func (uc *UserCase) Registration(dto *entity.LoginDto) (string, error) {
	token, err := uc.del.Registration(dto)
	return token, err
}
