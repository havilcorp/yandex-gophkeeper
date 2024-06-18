package usecase

import (
	"bytes"

	"yandex-gophkeeper-client/internal/entity"
	"yandex-gophkeeper-client/internal/storage"

	"github.com/sirupsen/logrus"
)

type UserCase struct {
	del   storage.Delivery
	local storage.LocalStorager
}

func New(del storage.Delivery, local storage.LocalStorager) *UserCase {
	return &UserCase{
		del:   del,
		local: local,
	}
}

func (uc *UserCase) SetToken(token string) {
	uc.del.SetToken(token)
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

func (uc *UserCase) Sync() error {
	listStore, err := uc.del.GetAll()
	if err != nil {
		logrus.Error(err)
		return err
	}
	listLocal, err := uc.local.GetAll()
	if err != nil {
		logrus.Error(err)
		return err
	}
	for _, localItem := range *listLocal {
		isFind := false
		for _, storeItem := range *listStore {
			if bytes.Equal(storeItem.Data, localItem.Data) {
				isFind = true
				break
			}
		}
		if !isFind {
			logrus.Info("Sync local to server")
			err := uc.del.Save(&localItem)
			if err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	for _, storeItem := range *listStore {
		isFind := false
		for _, localItem := range *listLocal {
			if bytes.Equal(localItem.Data, storeItem.Data) {
				isFind = true
				break
			}
		}
		if !isFind {
			logrus.Info("Sync server to local")
			err := uc.local.Save(&storeItem)
			if err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	return nil
}

func (uc *UserCase) GetByServerAll() (*[]entity.ItemDto, error) {
	return uc.del.GetAll()
}

func (uc *UserCase) GetByLocalAll() (*[]entity.ItemDto, error) {
	return uc.local.GetAll()
}
