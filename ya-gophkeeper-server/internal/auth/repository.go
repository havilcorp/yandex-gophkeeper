package auth

import "yandex-gophkeeper-server/internal/auth/entity"

// Repository интерфейс взаимодействия бизнес логики и репозитория
type Repository interface {
	GetUser(email string) (*entity.User, error)
	CreateUser(email string, hashPassword string) (*entity.User, error)
}
