package auth

import "ya-gophkeeper-server/internal/auth/entity"

type UserCase interface {
	Login(email string, password string) (*entity.User, error)
	Registration(email string, password string) (*entity.User, error)
}
