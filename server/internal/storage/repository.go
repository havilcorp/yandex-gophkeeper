package storage

import "ya-gophkeeper-server/internal/storage/entity"

type Repository interface {
	Save(*entity.CreateDto) error
	GetById(id int) (*entity.Item, error)
	GetAll() (*[]entity.Item, error)
	Remove(id int) error
}
