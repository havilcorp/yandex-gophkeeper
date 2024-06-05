package usecase

import (
	"ya-gophkeeper-client/internal/store"
	"ya-gophkeeper-client/internal/store/entity"

	"github.com/sirupsen/logrus"
)

type UserCase struct {
	del store.Delivery
}

func New(del store.Delivery) *UserCase {
	return &UserCase{
		del: del,
	}
}

func (uc *UserCase) Save(data []byte) error {
	err := uc.del.Save(&entity.SaveDto{
		Data: data,
		Meta: "test",
	})
	return err
}

func (uc *UserCase) GetById(id string) (*entity.SaveDto, error) {
	data, err := uc.del.GetByID(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return data, err
}
