package grpc

import (
	"context"
	"errors"
	"testing"

	"yandex-gophkeeper-server/internal/config"
	"yandex-gophkeeper-server/internal/storage/entity"
	"yandex-gophkeeper-server/internal/storage/mocks"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"
)

func Test_handler_Save(t *testing.T) {
	uc := mocks.NewUseCase(t)

	uc.On("Save", 1, &entity.CreateDto{
		Data: []byte(""),
		Meta: "",
	}).Return(nil)

	uc.On("Save", 1, &entity.CreateDto{
		Data: []byte(""),
		Meta: "error",
	}).Return(errors.New(""))

	type args struct {
		ctx context.Context
		in  *pb.SaveRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				ctx: context.WithValue(context.Background(), "X-User-ID", "1"),
				in: &pb.SaveRequest{
					Data: []byte(""),
					Meta: "",
				},
			},
			wantErr: false,
		},
		{
			name: "err user",
			args: args{
				ctx: context.WithValue(context.Background(), "X-User-ID", "notnum"),
				in: &pb.SaveRequest{
					Data: []byte(""),
					Meta: "",
				},
			},
			wantErr: true,
		},
		{
			name: "err user",
			args: args{
				ctx: context.WithValue(context.Background(), "X-User-ID", "1"),
				in: &pb.SaveRequest{
					Data: []byte(""),
					Meta: "error",
				},
			},
			wantErr: true,
		},
	}
	conf := config.Config{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&conf, uc)
			_, err := h.Save(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_handler_GetAll(t *testing.T) {
	uc := mocks.NewUseCase(t)

	uc.On("GetAll", 1).Return(&[]entity.Item{{
		ID:     1,
		UserId: 1,
		Data:   []byte(""),
		Meta:   "",
	}}, nil)
	uc.On("GetAll", 2).Return(nil, errors.New(""))

	type args struct {
		ctx context.Context
		in  *pb.GetAllRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				ctx: context.WithValue(context.Background(), "X-User-ID", "1"),
				in:  &pb.GetAllRequest{},
			},
			wantErr: false,
		},
		{
			name: "err user",
			args: args{
				ctx: context.WithValue(context.Background(), "X-User-ID", "notnum"),
				in:  &pb.GetAllRequest{},
			},
			wantErr: true,
		},
		{
			name: "err GetAll",
			args: args{
				ctx: context.WithValue(context.Background(), "X-User-ID", "2"),
				in:  &pb.GetAllRequest{},
			},
			wantErr: true,
		},
	}
	conf := config.Config{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&conf, uc)
			_, err := h.GetAll(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
