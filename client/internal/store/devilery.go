package store

import "ya-gophkeeper-client/internal/store/entity"

type Delivery interface {
	Save(dto *entity.SaveDto) error
	GetByID(id string) (*entity.SaveDto, error)
}
