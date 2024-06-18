package storage

import "yandex-gophkeeper-server/internal/storage/entity"

// Repository интерфейс взаимодействия бизнес логики и репозитория
type Repository interface {
	Save(userID int, dto *entity.CreateDto) error
	GetAll(userID int) (*[]entity.Item, error)
}
