// Package grpc транспортный уровень хранилища
package grpc

import (
	"context"
	"strconv"

	"yandex-gophkeeper-server/internal/config"
	"yandex-gophkeeper-server/internal/storage"
	"yandex-gophkeeper-server/internal/storage/entity"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"
)

type handler struct {
	pb.UnimplementedSaveServer
	conf *config.Config
	uc   storage.UseCase
}

// NewHandler получить экземпляр хендлера
func NewHandler(conf *config.Config, uc storage.UseCase) *handler {
	return &handler{
		conf: conf,
		uc:   uc,
	}
}

// Save сохранить данные
func (h *handler) Save(ctx context.Context, in *pb.SaveRequest) (*pb.SaveResponse, error) {
	var response pb.SaveResponse
	userStr := ctx.Value("X-User-ID").(string)
	user, err := strconv.Atoi(userStr)
	if err != nil {
		return &response, err
	}
	err = h.uc.Save(user, &entity.CreateDto{
		Data: in.Data,
		Meta: in.Meta,
	})
	if err != nil {
		return &response, err
	}
	return &response, nil
}

// GetAll получить данные
func (h *handler) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	var response pb.GetAllResponse
	userStr := ctx.Value("X-User-ID").(string)
	user, err := strconv.Atoi(userStr)
	if err != nil {
		return &response, err
	}
	items, err := h.uc.GetAll(user)
	if err != nil {
		return &response, err
	}
	for _, item := range *items {
		response.Items = append(response.Items, &pb.SaveRequest{
			Data: item.Data,
			Meta: item.Meta,
		})
	}
	return &response, nil
}
