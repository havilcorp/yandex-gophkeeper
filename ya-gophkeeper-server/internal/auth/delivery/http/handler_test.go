package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"yandex-gophkeeper-server/internal/auth/entity"
	"yandex-gophkeeper-server/internal/auth/mocks"
	"yandex-gophkeeper-server/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	r := chi.NewRouter()
	conf := &config.Config{}
	uc := mocks.NewUseCase(t)
	h := NewHandler(conf, uc)
	h.Register(r)
}

func Test_handler_login(t *testing.T) {
	uc := mocks.NewUseCase(t)
	conf := &config.Config{JWTKey: "jwt"}

	uc.On("Login", "mail@mail.ru", "toor").Return(&entity.User{ID: 1, Email: "", Password: ""}, nil)
	uc.On("Login", "test@mail.ru", "toor").Return(nil, entity.ErrUserNotFound)
	uc.On("Login", "mail@mail.ru", "not").Return(nil, entity.ErrWrongPassword)
	uc.On("Login", "aaa", "aaa").Return(nil, errors.New(""))

	type args struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good",
			args: args{
				statusCode: 200,
				body:       `{"email": "mail@mail.ru", "password": "toor"}`,
			},
		},
		{
			name: "bad json",
			args: args{
				statusCode: 400,
				body:       `{email: "mail@mail.ru", password: "toor"}`,
			},
		},
		{
			name: "user not found",
			args: args{
				statusCode: 400,
				body:       `{"email": "test@mail.ru", "password": "toor"}`,
			},
		},
		{
			name: "wrong password",
			args: args{
				statusCode: 401,
				body:       `{"email": "mail@mail.ru", "password": "not"}`,
			},
		},
		{
			name: "error",
			args: args{
				statusCode: 500,
				body:       `{"email": "aaa", "password": "aaa"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.args.body))
			h := NewHandler(conf, uc)
			h.login(w, r)
			res := w.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}

func Test_handler_registration(t *testing.T) {
	uc := mocks.NewUseCase(t)
	conf := &config.Config{JWTKey: "jwt"}
	conf.JWTKey = "jwt"

	uc.On("Registration", "mail@mail.ru", "toor").Return(&entity.User{ID: 1, Email: "", Password: ""}, nil)
	uc.On("Registration", "test@mail.ru", "toor").Return(nil, entity.ErrUserExists)
	uc.On("Registration", "aaa", "aaa").Return(nil, errors.New(""))

	type args struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good",
			args: args{
				statusCode: 200,
				body:       `{"email": "mail@mail.ru", "password": "toor"}`,
			},
		},
		{
			name: "bad json",
			args: args{
				statusCode: 400,
				body:       `{email: "mail@mail.ru", password: "toor"}`,
			},
		},
		{
			name: "user exists",
			args: args{
				statusCode: 400,
				body:       `{"email": "test@mail.ru", "password": "toor"}`,
			},
		},
		{
			name: "error",
			args: args{
				statusCode: 500,
				body:       `{"email": "aaa", "password": "aaa"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/registration", strings.NewReader(tt.args.body))
			h := NewHandler(conf, uc)
			h.registration(w, r)
			res := w.Result()
			assert.Equal(t, tt.args.statusCode, res.StatusCode)
		})
	}
}
