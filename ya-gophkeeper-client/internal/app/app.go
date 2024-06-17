package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"

	authDevivery "yandex-gophkeeper-client/internal/auth/delivery/http"
	authEntity "yandex-gophkeeper-client/internal/auth/entity"
	authUseCase "yandex-gophkeeper-client/internal/auth/usecase"
	"yandex-gophkeeper-client/internal/sqlite"

	storeDeviveryGrpc "yandex-gophkeeper-client/internal/store/delivery/grpc"
	storeUseCase "yandex-gophkeeper-client/internal/store/usecase"

	cmdInterface "yandex-gophkeeper-client/internal/cli"
	"yandex-gophkeeper-client/internal/config"

	"yandex-gophkeeper-client/pkg/crypto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Start() {
	conf := config.New()

	// TLS

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
	certConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}
	credTLS := credentials.NewTLS(certConfig)

	// Client GRPC

	connectGRPC, err := grpc.NewClient(conf.AddressGRPC, grpc.WithTransportCredentials(credTLS))
	if err != nil {
		logrus.Panic(err)
	}
	defer connectGRPC.Close()
	clientGRPC := pb.NewSaveClient(connectGRPC)

	// Client HTTP

	clientHTTP := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: certConfig,
		},
	}

	// SQLite

	db, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		logrus.Panic(err)
	}
	defer db.Close()
	sqlt := sqlite.New(db)
	if err := sqlt.Migration(); err != nil {
		logrus.Panic(err)
	}

	// Architecture

	authD := authDevivery.New(conf, &clientHTTP)
	authUC := authUseCase.New(authD)
	// storeDelHTTP := storeDeviveryHttp.New(conf, &clientHTTP)
	storeDelGRPC := storeDeviveryGrpc.New(conf, clientGRPC)
	storeUC := storeUseCase.New(storeDelGRPC, sqlt)

	var token string
	var cryptoKey []byte

	// CLI

	cli := cmdInterface.New()

	cli.Register("login", func() {
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
		token, err = authUC.Login(&authEntity.LoginDto{
			Email:    email,
			Password: pass,
		})
		if err != nil {
			logrus.Error(err)
			return
		}
		cli.Println("Вы успешно авторизовались!")
		storeDelGRPC.SetToken(token)
	})

	cli.Register("registration", func() {
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
		token, err = authUC.Registration(&authEntity.LoginDto{
			Email:    email,
			Password: pass,
		})
		if err != nil {
			logrus.Error(err)
			return
		}
		storeDelGRPC.SetToken(token)
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
		list, err := storeUC.GetByLocalAll()
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

	cli.Register("sync", func() {
		err := storeUC.Sync()
		if err != nil {
			logrus.Error(err)
		}
	})

	cli.Run()
}
