package auth

import "yandex-gophkeeper-server/internal/auth/entity"

// UseCase интерфейс взаимодействия транспортного уровня и бизнес логики
type UseCase interface {
	Login(email string, password string) (*entity.User, error)
	Registration(email string, password string) (*entity.User, error)
}
