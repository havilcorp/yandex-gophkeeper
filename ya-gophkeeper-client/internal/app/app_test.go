package app

import (
	"testing"

	authEntity "yandex-gophkeeper-client/internal/auth/entity"
	authMock "yandex-gophkeeper-client/internal/auth/mocks"
	cliMock "yandex-gophkeeper-client/internal/cli/mocks"
	storageMock "yandex-gophkeeper-client/internal/storage/mocks"

	"github.com/stretchr/testify/mock"
)

func Test_getCert(t *testing.T) {
	type args struct {
		ca  string
		crt string
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				ca:  "../../tls/ca.crt",
				crt: "../../tls/server.crt",
				key: "../../tls/server.key",
			},
			wantErr: false,
		},
		{
			name: "err1",
			args: args{
				ca:  "../../tls/notfound.crt",
				crt: "../../tls/server.crt",
				key: "../../tls/server.key",
			},
			wantErr: true,
		},
		{
			name: "err2",
			args: args{
				ca:  "../../tls/ca.crt",
				crt: "../../tls/notfound.crt",
				key: "../../tls/server.key",
			},
			wantErr: true,
		},
		{
			name: "err3",
			args: args{
				ca:  "../../tls/ca.crt",
				crt: "../../tls/server.crt",
				key: "../../tls/notfound.key",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getCert(tt.args.ca, tt.args.crt, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_cliStart(t *testing.T) {
	cli := cliMock.NewCLIer(t)
	authUC := authMock.NewUseCase(t)
	storageUC := storageMock.NewUseCase(t)

	cli.On("Register", "login", mock.Anything).Return()
	cli.On("Register", "registration", mock.Anything).Return()
	cli.On("Register", "crypto", mock.Anything).Return()
	cli.On("Register", "add", mock.Anything).Return()
	cli.On("Register", "list", mock.Anything).Return()
	cli.On("Register", "sync", mock.Anything).Return()

	cli.On("GetUserPrint", "Почта: ").Return("mail@mail.ru", nil)
	cli.On("GetHideUserPrint", "Пароль: ").Return("password", nil)
	authUC.On("Login", &authEntity.LoginDto{
		Email:    "mail@mail.ru",
		Password: "password",
	}).Return("token", nil)
	cli.On("Println", "Вы успешно авторизовались!").Return()
	storageUC.On("SetToken", "token").Return()

	mc := NewMainController(cli, authUC, storageUC)

	mc.Register()

	mc.login()
}
