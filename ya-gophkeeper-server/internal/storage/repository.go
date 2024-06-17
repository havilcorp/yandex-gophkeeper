package storage

import "yandex-gophkeeper-server/internal/storage/entity"

type Repository interface {
	Save(userID int, dto *entity.CreateDto) error
	GetAll(userID int) (*[]entity.Item, error)
}
