package storage

import "ya-gophkeeper-server/internal/storage/entity"

type UserCase interface {
	Save(userId int, dto *entity.CreateDto) error
	GetById(id int) (*entity.Item, error)
	GetAll(userID int) (*[]entity.Item, error)
	Remove(id int) error
}
