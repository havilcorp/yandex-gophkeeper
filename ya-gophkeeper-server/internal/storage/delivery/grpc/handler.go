package grpc

import (
	"context"

	"yandex-gophkeeper-server/internal/config"
	"yandex-gophkeeper-server/internal/storage"
	"yandex-gophkeeper-server/internal/storage/entity"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"
	"github.com/sirupsen/logrus"
)

type handler struct {
	pb.UnimplementedSaveServer
	conf *config.Config
	uc   storage.UserCase
}

func NewHandler(conf *config.Config, uc storage.UserCase) *handler {
	return &handler{
		conf: conf,
		uc:   uc,
	}
}

func (h *handler) Save(ctx context.Context, in *pb.SaveRequest) (*pb.SaveResponse, error) {
	var response pb.SaveResponse
	err := h.uc.Save(1, &entity.CreateDto{
		Data: in.Data,
		Meta: in.Meta,
	})
	if err != nil {
		logrus.Error(err)
		return &response, err
	}
	return &response, nil
}

func (h *handler) GetAll(ctx context.Context, in *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	var response pb.GetAllResponse
	items, err := h.uc.GetAll(1)
	if err != nil {
		logrus.Error(err)
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
