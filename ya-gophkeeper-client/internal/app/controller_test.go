package app

import (
	"testing"

	authEntity "yandex-gophkeeper-client/internal/auth/entity"
	authMock "yandex-gophkeeper-client/internal/auth/mocks"
	cliMock "yandex-gophkeeper-client/internal/cli/mocks"
	"yandex-gophkeeper-client/internal/entity"
	storageMock "yandex-gophkeeper-client/internal/storage/mocks"
	"yandex-gophkeeper-client/pkg/crypto"

	"github.com/stretchr/testify/mock"
)

func TestMainController_Register(t *testing.T) {
	cli := cliMock.NewCLIer(t)
	cli.On("Register", "login", mock.Anything).Return()
	cli.On("Register", "registration", mock.Anything).Return()
	cli.On("Register", "crypto", mock.Anything).Return()
	cli.On("Register", "add", mock.Anything).Return()
	cli.On("Register", "list", mock.Anything).Return()
	cli.On("Register", "sync", mock.Anything).Return()
	mc := NewMainController(cli, nil, nil)
	mc.Register()
}

func TestMainController_login(t *testing.T) {
	cli := cliMock.NewCLIer(t)
	authUC := authMock.NewUseCase(t)
	storageUC := storageMock.NewUseCase(t)

	t.Run("good", func(t *testing.T) {
		cli.On("GetUserPrint", "Почта: ").Return("mail@mail.ru", nil)
		cli.On("GetHideUserPrint", "Пароль: ").Return("password", nil)
		authUC.On("Login", &authEntity.LoginDto{
			Email:    "mail@mail.ru",
			Password: "password",
		}).Return("token", nil)
		cli.On("Println", "Вы успешно авторизовались!").Return()
		storageUC.On("SetToken", "token").Return()
		mc := NewMainController(cli, authUC, storageUC)
		mc.login()
	})
}

func TestMainController_registration(t *testing.T) {
	cli := cliMock.NewCLIer(t)
	authUC := authMock.NewUseCase(t)
	storageUC := storageMock.NewUseCase(t)

	t.Run("good", func(t *testing.T) {
		cli.On("GetUserPrint", "Почта: ").Return("mail@mail.ru", nil)
		cli.On("GetHideUserPrint", "Пароль: ").Return("password", nil)
		cli.On("GetHideUserPrint", "Повторите пароль: ").Return("password", nil)
		// cli.On("Println", "Пароли не совпадают!").Return()

		authUC.On("Registration", &authEntity.LoginDto{
			Email:    "mail@mail.ru",
			Password: "password",
		}).Return("token", nil)
		storageUC.On("SetToken", "token").Return()
		mc := NewMainController(cli, authUC, storageUC)
		mc.registration()
	})
}

func TestMainController_crypto(t *testing.T) {
	cli := cliMock.NewCLIer(t)
	authUC := authMock.NewUseCase(t)
	storageUC := storageMock.NewUseCase(t)

	t.Run("good", func(t *testing.T) {
		cli.On("GetHideUserPrint", "Введите пароль для шифрования данных: ").Return("pass", nil)
		mc := NewMainController(cli, authUC, storageUC)
		mc.crypto()
	})
}

func TestMainController_add(t *testing.T) {
	cryptoKey, err := crypto.GetKey([]byte("pass"), []byte("salt"))
	if err != nil {
		t.Error(err)
	}

	t.Run("add login pass", func(t *testing.T) {
		cli := cliMock.NewCLIer(t)
		authUC := authMock.NewUseCase(t)
		storageUC := storageMock.NewUseCase(t)

		cli.On("Println", "1 - Добавить логин:пароль")
		cli.On("Println", "2 - Добавить данные банковской карты")
		cli.On("Println", "3 - Добавить произвольный текст")

		cli.On("GetUserPrint", "login: ").Return("admin", nil)
		cli.On("GetHideUserPrint", "pass: ").Return("pass", nil)

		cli.On("GetUserPrint", "Введите число: ").Return("1", nil)

		mc := NewMainController(cli, authUC, storageUC)
		mc.cryptoKey = cryptoKey
		storageUC.On("Save", mock.Anything, `{"type":"loginpass"}`).Return(nil)

		mc.add()
	})

	t.Run("add bank card", func(t *testing.T) {
		cli := cliMock.NewCLIer(t)
		authUC := authMock.NewUseCase(t)
		storageUC := storageMock.NewUseCase(t)

		cli.On("Println", "1 - Добавить логин:пароль")
		cli.On("Println", "2 - Добавить данные банковской карты")
		cli.On("Println", "3 - Добавить произвольный текст")

		cli.On("GetUserPrint", "number: ").Return("2231 2221 3455 3211", nil)
		cli.On("GetUserPrint", "date: ").Return("02/29", nil)
		cli.On("GetHideUserPrint", "cvv: ").Return("123", nil)

		cli.On("GetUserPrint", "Введите число: ").Return("2", nil)

		mc := NewMainController(cli, authUC, storageUC)
		mc.cryptoKey = cryptoKey
		storageUC.On("Save", mock.Anything, `{"type":"bankcard"}`).Return(nil)

		mc.add()
	})

	t.Run("add bank card", func(t *testing.T) {
		cli := cliMock.NewCLIer(t)
		authUC := authMock.NewUseCase(t)
		storageUC := storageMock.NewUseCase(t)

		cli.On("Println", "1 - Добавить логин:пароль")
		cli.On("Println", "2 - Добавить данные банковской карты")
		cli.On("Println", "3 - Добавить произвольный текст")

		cli.On("GetUserPrint", "text: ").Return("Hello World", nil)

		cli.On("GetUserPrint", "Введите число: ").Return("3", nil)

		mc := NewMainController(cli, authUC, storageUC)
		mc.cryptoKey = cryptoKey
		storageUC.On("Save", mock.Anything, `{"type":"text"}`).Return(nil)

		mc.add()
	})
}

func TestMainController_list(t *testing.T) {
	cli := cliMock.NewCLIer(t)
	authUC := authMock.NewUseCase(t)
	storageUC := storageMock.NewUseCase(t)

	cli.On("Println", mock.Anything).Return()

	mc := NewMainController(cli, authUC, storageUC)
	cryptoKey, err := crypto.GetKey([]byte("pass"), []byte("salt"))
	if err != nil {
		t.Error(err)
	}
	mc.cryptoKey = cryptoKey

	dataEncrypt, err := crypto.Encrypt(mc.cryptoKey, []byte("Hello World"))
	if err != nil {
		t.Error(err)
	}

	storageUC.On("GetByLocalAll").Return(&[]entity.ItemDto{
		{
			Data: dataEncrypt,
			Meta: `{"type":"text"}`,
		},
	}, nil)

	mc.list()
}

func TestMainController_sync(t *testing.T) {
	cli := cliMock.NewCLIer(t)
	authUC := authMock.NewUseCase(t)
	storageUC := storageMock.NewUseCase(t)

	storageUC.On("Sync").Return(nil)

	mc := NewMainController(cli, authUC, storageUC)
	mc.sync()
}
