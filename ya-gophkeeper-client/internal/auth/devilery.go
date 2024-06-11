package auth

import "ya-gophkeeper-client/internal/auth/entity"

type Delivery interface {
	Login(dto *entity.LoginDto) (string, error)
	Registration(dto *entity.LoginDto) (string, error)
}
