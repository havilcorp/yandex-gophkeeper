package usecase

import (
	"yandex-gophkeeper-server/internal/storage"
	"yandex-gophkeeper-server/internal/storage/entity"
)

type usecace struct {
	repo storage.Repository
}

func New(repo storage.Repository) *usecace {
	return &usecace{
		repo: repo,
	}
}

// TODO: add logic
func (uc *usecace) Save(userID int, dto *entity.CreateDto) error {
	return uc.repo.Save(userID, dto)
}

func (uc *usecace) GetById(id int) (*entity.Item, error) {
	return uc.repo.GetById(id)
}

func (uc *usecace) GetAll(userID int) (*[]entity.Item, error) {
	return uc.repo.GetAll(userID)
}

func (uc *usecace) Remove(id int) error {
	return uc.repo.Remove(id)
}
