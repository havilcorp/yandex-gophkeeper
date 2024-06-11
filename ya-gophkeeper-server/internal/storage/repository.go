package storage

import "ya-gophkeeper-server/internal/storage/entity"

type Repository interface {
	Save(userID int, dto *entity.CreateDto) error
	GetById(id int) (*entity.Item, error)
	GetAll(userID int) (*[]entity.Item, error)
	Remove(id int) error
}
