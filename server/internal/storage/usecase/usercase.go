package usecase

import (
	"ya-gophkeeper-server/internal/storage"
	"ya-gophkeeper-server/internal/storage/entity"
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
func (uc *usecace) Save(userId int, dto *entity.CreateDto) error {
	return uc.repo.Save(dto)
}

func (uc *usecace) GetById(id int) (*entity.Item, error) {
	return uc.repo.GetById(id)
}

func (uc *usecace) GetAll() (*[]entity.Item, error) {
	return uc.repo.GetAll()
}

func (uc *usecace) Remove(id int) error {
	return uc.repo.Remove(id)
}
