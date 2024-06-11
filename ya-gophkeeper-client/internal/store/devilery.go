package store

import "ya-gophkeeper-client/internal/entity"

type Delivery interface {
	Save(dto *entity.ItemDto) error
	GetByID(id string) (*entity.ItemDto, error)
	GetList() (*[]entity.ItemDto, error)
}
