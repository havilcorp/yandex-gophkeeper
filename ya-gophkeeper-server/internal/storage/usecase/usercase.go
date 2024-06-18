// Package usecase пакет бизнес логики хранилища
package usecase

import (
	"yandex-gophkeeper-server/internal/storage"
	"yandex-gophkeeper-server/internal/storage/entity"
)

type usecace struct {
	repo storage.Repository
}

// New получить экземпляр структуры
func New(repo storage.Repository) *usecace {
	return &usecace{
		repo: repo,
	}
}

// Save сохранить данные
func (uc *usecace) Save(userID int, dto *entity.CreateDto) error {
	return uc.repo.Save(userID, dto)
}

// GetAll получить все данные
func (uc *usecace) GetAll(userID int) (*[]entity.Item, error) {
	return uc.repo.GetAll(userID)
}
