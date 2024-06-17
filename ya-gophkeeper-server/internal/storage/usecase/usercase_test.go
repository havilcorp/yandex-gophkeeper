package usecase

import (
	"testing"

	"yandex-gophkeeper-server/internal/storage/entity"
	"yandex-gophkeeper-server/internal/storage/mocks"
)

func Test_usecace_Save(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.On("Save", 1, &entity.CreateDto{}).Return(nil)
	uc := New(repo)
	err := uc.Save(1, &entity.CreateDto{})
	if err != nil {
		t.Error(err)
	}
}

func Test_usecace_GetAll(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.On("GetAll", 1).Return(nil, nil)
	uc := New(repo)
	_, err := uc.GetAll(1)
	if err != nil {
		t.Error(err)
	}
}
