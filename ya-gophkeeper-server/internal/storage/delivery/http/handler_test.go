package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"yandex-gophkeeper-server/internal/config"
	"yandex-gophkeeper-server/internal/storage/entity"
	"yandex-gophkeeper-server/internal/storage/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	r := chi.NewRouter()
	conf := &config.Config{JWTKey: "jwt"}
	uc := mocks.NewUseCase(t)
	h := NewHandler(conf, uc)
	h.Register(r)
}

func Test_handler_Save(t *testing.T) {
	uc := mocks.NewUseCase(t)

	uc.On("Save", 1, &entity.CreateDto{
		Data: []byte(""),
		Meta: "",
	}).Return(nil)

	uc.On("Save", 1, &entity.CreateDto{
		Data: []byte(""),
		Meta: "err",
	}).Return(errors.New(""))

	type args struct {
		statusCode int
		body       string
		userID     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good",
			args: args{
				statusCode: 200,
				body:       `{"data": "", "meta": ""}`,
				userID:     "1",
			},
		},
		{
			name: "error header token",
			args: args{
				statusCode: 401,
				body:       `{"data": "", "meta": ""}`,
				userID:     "",
			},
		},
		{
			name: "error user id",
			args: args{
				statusCode: 500,
				body:       `{"data": "", "meta": ""}`,
				userID:     "err",
			},
		},
		{
			name: "error save",
			args: args{
				statusCode: 500,
				body:       `{"data": "", "meta": "err"}`,
				userID:     "1",
			},
		},
	}
	conf := config.Config{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/storage/", strings.NewReader(tt.args.body))
			r.Header.Add("X-User-ID", tt.args.userID)
			h := NewHandler(&conf, uc)
			h.save(w, r)
			res := w.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}

func Test_handler_getAll(t *testing.T) {
	uc := mocks.NewUseCase(t)

	uc.On("GetAll", 1).Return(&[]entity.Item{{
		ID:     1,
		UserId: 1,
		Data:   []byte(""),
		Meta:   "",
	}}, nil)
	uc.On("GetAll", 2).Return(nil, errors.New(""))

	type args struct {
		statusCode int
		body       string
		userID     string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good",
			args: args{
				statusCode: 200,
				userID:     "1",
			},
		},
		{
			name: "error header token",
			args: args{
				statusCode: 401,
				userID:     "",
			},
		},
		{
			name: "error user id",
			args: args{
				statusCode: 500,
				userID:     "err",
			},
		},
		{
			name: "error save",
			args: args{
				statusCode: 500,
				userID:     "2",
			},
		},
	}
	conf := config.Config{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/storage/", strings.NewReader(tt.args.body))
			r.Header.Add("X-User-ID", tt.args.userID)
			h := NewHandler(&conf, uc)
			h.getAll(w, r)
			res := w.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}
