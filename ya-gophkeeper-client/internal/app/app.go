package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	authDevivery "ya-gophkeeper-client/internal/auth/delivery/http"
	authUseCase "ya-gophkeeper-client/internal/auth/usecase"
	"ya-gophkeeper-client/internal/sqlite"

	storeDevivery "ya-gophkeeper-client/internal/store/delivery/http"
	storeUseCase "ya-gophkeeper-client/internal/store/usecase"

	"ya-gophkeeper-client/internal/auth/entity"
	cmdInterface "ya-gophkeeper-client/internal/cli"
	"ya-gophkeeper-client/internal/config"

	"ya-gophkeeper-client/pkg/crypto"

	"github.com/sirupsen/logrus"
)

func Start() {
	conf := config.New()

	caCert, err := os.ReadFile("./tls/ca.crt")
	if err != nil {
		logrus.Error(err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		logrus.Info("Not ok")
	}
	cert, err := tls.LoadX509KeyPair("./tls/server.crt", "./tls/server.key")
	if err != nil {
		logrus.Error(err)
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	db, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		logrus.Panic(err)
	}
	defer db.Close()

	sqlt := sqlite.New(db)
	if err := sqlt.Migration(); err != nil {
		logrus.Panic(err)
	}

	authD := authDevivery.New(conf, &client)
	authUC := authUseCase.New(authD)
	storeD := storeDevivery.New(conf, &client)
	storeUC := storeUseCase.New(storeD, sqlt)

	var token string
	var cryptoKey []byte

	cli := cmdInterface.New()

	cli.Register("auth", func() {
		email, err := cli.GetUserPrint("Почта: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		pass, err := cli.GetHideUserPrint("Пароль: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		token, err = authUC.Login(&entity.LoginDto{
			Email:    email,
			Password: pass,
		})
		if err != nil {
			logrus.Error(err)
			return
		}
		cli.Println("Вы успешно авторизовались!")
		storeD.SetToken(token)
	})

	cli.Register("reg", func() {
		email, err := cli.GetUserPrint("Почта: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		pass, err := cli.GetHideUserPrint("Пароль: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		repass, err := cli.GetHideUserPrint("Повторите пароль: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		if pass != repass {
			cli.Println("Пароли не совпадают!")
			return
		}
		token, err = authUC.Registration(&entity.LoginDto{
			Email:    email,
			Password: pass,
		})
		if err != nil {
			logrus.Error(err)
			return
		}
		storeD.SetToken(token)
	})

	cli.Register("crypto", func() {
		cryptoPass, err := cli.GetHideUserPrint("Введите пароль для шифрования данных: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		cryptoKey, err = crypto.GetKey([]byte(cryptoPass), []byte("salt"))
		if err != nil {
			logrus.Error(err)
			return
		}
	})

	cli.Register("add", func() {
		if len(cryptoKey) == 0 {
			cli.Call("crypto")
		}
		cli.Println("1 - Добавить логин:пароль")
		cli.Println("2 - Добавить данные банковской карты")
		cli.Println("3 - Добавить произвольный текст")
		num, err := cli.GetUserPrint("Введите число: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		var data string
		var meta string
		switch num {
		case "1":
			{
				login, err := cli.GetUserPrint("login: ")
				if err != nil {
					logrus.Error(err)
					return
				}
				pass, err := cli.GetHideUserPrint("pass: ")
				if err != nil {
					logrus.Error(err)
					return
				}
				data = fmt.Sprintf(`{"login": "%s", "pass": "%s"}`, login, pass)
				meta = `{"type":"loginpass"}`
			}
		case "2":
			{
				number, err := cli.GetUserPrint("number: ")
				if err != nil {
					logrus.Error(err)
					return
				}
				date, err := cli.GetUserPrint("date: ")
				if err != nil {
					logrus.Error(err)
					return
				}
				cvv, err := cli.GetHideUserPrint("cvv: ")
				if err != nil {
					logrus.Error(err)
					return
				}
				data = fmt.Sprintf(`{"number": "%s", "date": "%s", "cvv": "%s"}`, number, date, cvv)
				meta = `{"type":"bankcard"}`
			}
		case "3":
			{
				data, err = cli.GetUserPrint("text: ")
				if err != nil {
					logrus.Error(err)
					return
				}
				meta = `{"type":"text"}`
			}
		}
		dataEncrypt, err := crypto.Encrypt(cryptoKey, []byte(data))
		if err != nil {
			logrus.Error(err)
			return
		}
		err = storeUC.Save(dataEncrypt, meta)
		if err != nil {
			logrus.Error(err)
			return
		}
	})

	cli.Register("list", func() {
		if len(cryptoKey) == 0 {
			cli.Call("crypto")
		}
		list, err := storeUC.GetList()
		if err != nil {
			logrus.Error(err)
			return
		}
		cli.Println("==_Your_secrets_============")
		for _, item := range *list {
			var meta struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal([]byte(item.Meta), &meta); err != nil {
				cli.Println("Не удалось рассшифровать")
				break
			}
			dataDecrypt, err := crypto.Decrypt(cryptoKey, item.Data)
			if err != nil {
				cli.Println("Не удалось рассшифровать")
			} else {
				cli.Println(fmt.Sprintf("- (%s): %s", meta.Type, dataDecrypt))
			}
		}
		cli.Println("============================")
	})

	cli.Run()
}
