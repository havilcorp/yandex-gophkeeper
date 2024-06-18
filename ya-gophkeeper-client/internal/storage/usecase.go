package storage

import "yandex-gophkeeper-client/internal/entity"

type UseCase interface {
	Save(data []byte, meta string) error
	Sync() error
	GetByServerAll() (*[]entity.ItemDto, error)
	GetByLocalAll() (*[]entity.ItemDto, error)
	SetToken(token string)
}
