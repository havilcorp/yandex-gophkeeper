package store

import "yandex-gophkeeper-client/internal/entity"

type LocalStorager interface {
	Save(item *entity.ItemDto) error
	GetAll() (*[]entity.ItemDto, error)
}
