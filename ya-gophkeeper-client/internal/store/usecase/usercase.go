package usecase

import (
	"ya-gophkeeper-client/internal/entity"
	"ya-gophkeeper-client/internal/store"

	"github.com/sirupsen/logrus"
)

type UserCase struct {
	del   store.Delivery
	local store.LocalStorager
}

func New(del store.Delivery, local store.LocalStorager) *UserCase {
	return &UserCase{
		del:   del,
		local: local,
	}
}

func (uc *UserCase) Save(data []byte, meta string) error {
	saveDto := entity.ItemDto{
		Data: data,
		Meta: meta,
	}
	err := uc.local.Save(&saveDto)
	// err = uc.del.Save(&entity.SaveDto{
	// 	Data: data,
	// 	Meta: meta,
	// })
	return err
}

func (uc *UserCase) GetById(id string) (*entity.ItemDto, error) {
	data, err := uc.del.GetByID(id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return data, err
}

func (uc *UserCase) GetList2() (*[]entity.ItemDto, error) {
	data, err := uc.del.GetList()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return data, err
}

func (uc *UserCase) GetList() (*[]entity.ItemDto, error) {
	return uc.local.GetAll()
}
