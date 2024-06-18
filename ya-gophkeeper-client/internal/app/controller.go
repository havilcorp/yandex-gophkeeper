package app

import (
	"encoding/json"
	"fmt"

	auth "yandex-gophkeeper-client/internal/auth"
	authEntity "yandex-gophkeeper-client/internal/auth/entity"
	"yandex-gophkeeper-client/internal/cli"
	storage "yandex-gophkeeper-client/internal/storage"
	"yandex-gophkeeper-client/pkg/crypto"

	"github.com/sirupsen/logrus"
)

type MainController struct {
	cli       cli.CLIer
	authUC    auth.UseCase
	storeUC   storage.UseCase
	cryptoKey []byte
}

// NewMainController контроллер для работы с командами
func NewMainController(cli cli.CLIer, authUC auth.UseCase, storeUC storage.UseCase) *MainController {
	return &MainController{
		cli:     cli,
		authUC:  authUC,
		storeUC: storeUC,
	}
}

// Register регистрация комманд
func (mc *MainController) Register() {
	mc.cli.Register("login", mc.login)
	mc.cli.Register("registration", mc.registration)
	mc.cli.Register("crypto", mc.crypto)
	mc.cli.Register("add", mc.add)
	mc.cli.Register("list", mc.list)
	mc.cli.Register("sync", mc.sync)
}

func (mc *MainController) login() {
	email, err := mc.cli.GetUserPrint("Почта: ")
	if err != nil {
		logrus.Error(err)
		return
	}
	pass, err := mc.cli.GetHideUserPrint("Пароль: ")
	if err != nil {
		logrus.Error(err)
		return
	}
	token, err := mc.authUC.Login(&authEntity.LoginDto{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	mc.cli.Println("Вы успешно авторизовались!")
	mc.storeUC.SetToken(token)
}

func (mc *MainController) registration() {
	email, err := mc.cli.GetUserPrint("Почта: ")
	if err != nil {
		logrus.Error(err)
		return
	}
	pass, err := mc.cli.GetHideUserPrint("Пароль: ")
	if err != nil {
		logrus.Error(err)
		return
	}
	repass, err := mc.cli.GetHideUserPrint("Повторите пароль: ")
	if err != nil {
		logrus.Error(err)
		return
	}
	if pass != repass {
		mc.cli.Println("Пароли не совпадают!")
		return
	}
	token, err := mc.authUC.Registration(&authEntity.LoginDto{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	mc.storeUC.SetToken(token)
}

func (mc *MainController) crypto() {
	cryptoPass, err := mc.cli.GetHideUserPrint("Введите пароль для шифрования данных: ")
	if err != nil {
		logrus.Error(err)
		return
	}
	mc.cryptoKey, err = crypto.GetKey([]byte(cryptoPass), []byte("salt"))
	if err != nil {
		logrus.Error(err)
		return
	}
}

func (mc *MainController) add() {
	if len(mc.cryptoKey) == 0 {
		mc.cli.CallFn("crypto")
	}
	mc.cli.Println("1 - Добавить логин:пароль")
	mc.cli.Println("2 - Добавить данные банковской карты")
	mc.cli.Println("3 - Добавить произвольный текст")
	num, err := mc.cli.GetUserPrint("Введите число: ")
	if err != nil {
		logrus.Error(err)
		return
	}
	var data string
	var meta string
	switch num {
	case "1":
		{
			login, err := mc.cli.GetUserPrint("login: ")
			if err != nil {
				logrus.Error(err)
				return
			}
			pass, err := mc.cli.GetHideUserPrint("pass: ")
			if err != nil {
				logrus.Error(err)
				return
			}
			data = fmt.Sprintf(`{"login": "%s", "pass": "%s"}`, login, pass)
			meta = `{"type":"loginpass"}`
		}
	case "2":
		{
			number, err := mc.cli.GetUserPrint("number: ")
			if err != nil {
				logrus.Error(err)
				return
			}
			date, err := mc.cli.GetUserPrint("date: ")
			if err != nil {
				logrus.Error(err)
				return
			}
			cvv, err := mc.cli.GetHideUserPrint("cvv: ")
			if err != nil {
				logrus.Error(err)
				return
			}
			data = fmt.Sprintf(`{"number": "%s", "date": "%s", "cvv": "%s"}`, number, date, cvv)
			meta = `{"type":"bankcard"}`
		}
	case "3":
		{
			data, err = mc.cli.GetUserPrint("text: ")
			if err != nil {
				logrus.Error(err)
				return
			}
			meta = `{"type":"text"}`
		}
	}
	dataEncrypt, err := crypto.Encrypt(mc.cryptoKey, []byte(data))
	if err != nil {
		logrus.Error(err)
		return
	}
	err = mc.storeUC.Save(dataEncrypt, meta)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func (mc *MainController) list() {
	if len(mc.cryptoKey) == 0 {
		mc.cli.CallFn("crypto")
	}
	list, err := mc.storeUC.GetByLocalAll()
	if err != nil {
		logrus.Error(err)
		return
	}
	mc.cli.Println("==_Your_secrets_============")
	for _, item := range *list {
		var meta struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(item.Meta), &meta); err != nil {
			mc.cli.Println("Не удалось рассшифровать")
			break
		}
		dataDecrypt, err := crypto.Decrypt(mc.cryptoKey, item.Data)
		if err != nil {
			mc.cli.Println("Не удалось рассшифровать")
		} else {
			mc.cli.Println(fmt.Sprintf("- (%s): %s", meta.Type, dataDecrypt))
		}
	}
	mc.cli.Println("============================")
}

func (mc *MainController) sync() {
	err := mc.storeUC.Sync()
	if err != nil {
		logrus.Error(err)
	}
}
