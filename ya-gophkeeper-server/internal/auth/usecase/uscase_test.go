package usecase

import (
	"errors"
	"testing"

	"yandex-gophkeeper-server/internal/auth/entity"
	"yandex-gophkeeper-server/internal/auth/mocks"
)

func Test_usecace_Login(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.On("GetUser", "mail@mail.ru").Return(&entity.User{Password: "1a1dc91c907325c69271ddf0c944bc72"}, nil)
	repo.On("GetUser", "notfound@mail.ru").Return(nil, entity.ErrUserNotFound)

	type args struct {
		err      error
		email    string
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good",
			args: args{
				err:      nil,
				email:    "mail@mail.ru",
				password: "pass",
			},
		},
		{
			name: "wrong password",
			args: args{
				err:      entity.ErrWrongPassword,
				email:    "mail@mail.ru",
				password: "wrong",
			},
		},
		{
			name: "user not found",
			args: args{
				err:      entity.ErrUserNotFound,
				email:    "notfound@mail.ru",
				password: "wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(repo)
			_, err := uc.Login(tt.args.email, tt.args.password)
			if !errors.Is(err, tt.args.err) {
				t.Error(err)
			}
		})
	}
}

func Test_usecace_Registration(t *testing.T) {
	repo := mocks.NewRepository(t)
	err := errors.New("ERROR")

	repo.On("GetUser", "mail@mail.ru").Return(nil, entity.ErrUserExists)
	repo.On("GetUser", "exists@mail.ru").Return(&entity.User{}, nil)
	repo.On("GetUser", "error@mail.ru").Return(nil, entity.ErrUserExists)

	repo.On("CreateUser", "mail@mail.ru", "1a1dc91c907325c69271ddf0c944bc72").Return(&entity.User{}, nil)
	repo.On("CreateUser", "error@mail.ru", "1a1dc91c907325c69271ddf0c944bc72").Return(nil, err)

	type args struct {
		err      error
		email    string
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good",
			args: args{
				err:      nil,
				email:    "mail@mail.ru",
				password: "pass",
			},
		},
		{
			name: "err user exists",
			args: args{
				err:      entity.ErrUserExists,
				email:    "exists@mail.ru",
				password: "",
			},
		},
		{
			name: "error create user",
			args: args{
				err:      err,
				email:    "error@mail.ru",
				password: "pass",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(repo)
			_, err := uc.Registration(tt.args.email, tt.args.password)
			if !errors.Is(err, tt.args.err) {
				t.Error(err)
			}
		})
	}
}
