package auth

import "ya-gophkeeper-server/internal/auth/entity"

type Repository interface {
	GetUser(email string) (*entity.User, error)
	CreateUser(email string, hashPassword string) (*entity.User, error)
}
