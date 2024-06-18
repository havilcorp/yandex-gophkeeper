// Package пакет авторизации
package auth

import "yandex-gophkeeper-client/internal/auth/entity"

type UseCase interface {
	Login(dto *entity.LoginDto) (string, error)
	Registration(dto *entity.LoginDto) (string, error)
}
