package storage

import "yandex-gophkeeper-client/internal/entity"

type Delivery interface {
	Save(dto *entity.ItemDto) error
	GetAll() (*[]entity.ItemDto, error)
	SetToken(token string)
}
