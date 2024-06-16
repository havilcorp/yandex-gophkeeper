package grpc

import (
	"context"
	"fmt"

	"yandex-gophkeeper-client/internal/config"
	"yandex-gophkeeper-client/internal/entity"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/metadata"
)

type handler struct {
	conf   *config.Config
	client pb.SaveClient
	token  string
}

func New(conf *config.Config, client pb.SaveClient) *handler {
	return &handler{
		conf:   conf,
		client: client,
	}
}

func (h *handler) SetToken(token string) {
	h.token = token
}

func (h *handler) Save(dto *entity.ItemDto) error {
	header := metadata.New(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", h.token),
	})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	_, err := h.client.Save(ctx, &pb.SaveRequest{
		Data: dto.Data,
		Meta: dto.Meta,
	})
	return err
}

func (h *handler) GetAll() (*[]entity.ItemDto, error) {
	header := metadata.New(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", h.token),
	})
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	data, err := h.client.GetAll(ctx, &pb.GetAllRequest{})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	items := make([]entity.ItemDto, 0)
	for _, item := range data.Items {
		items = append(items, entity.ItemDto{
			Data: item.Data,
			Meta: item.Meta,
		})
	}
	return &items, nil
}
