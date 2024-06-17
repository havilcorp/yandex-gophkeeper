package auth

import "yandex-gophkeeper-server/internal/auth/entity"

type UseCase interface {
	Login(email string, password string) (*entity.User, error)
	Registration(email string, password string) (*entity.User, error)
}
