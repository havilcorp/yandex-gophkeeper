package storage

import "yandex-gophkeeper-server/internal/storage/entity"

// UseCase интерфейс взаимодействия транспортного уровня и бизнес логики
type UseCase interface {
	Save(userId int, dto *entity.CreateDto) error
	GetAll(userID int) (*[]entity.Item, error)
}
