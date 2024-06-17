package storage

import "yandex-gophkeeper-server/internal/storage/entity"

type UseCase interface {
	Save(userId int, dto *entity.CreateDto) error
	GetAll(userID int) (*[]entity.Item, error)
}
