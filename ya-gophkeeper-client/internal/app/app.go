// Package app пакет для запуск клиента
package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"net/http"
	"os"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"

	authDevivery "yandex-gophkeeper-client/internal/auth/delivery/http"
	authUseCase "yandex-gophkeeper-client/internal/auth/usecase"
	"yandex-gophkeeper-client/internal/sqlite"

	storeDeviveryGrpc "yandex-gophkeeper-client/internal/storage/delivery/grpc"
	storeUseCase "yandex-gophkeeper-client/internal/storage/usecase"

	cmdInterface "yandex-gophkeeper-client/internal/cli"
	"yandex-gophkeeper-client/internal/config"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Start запускает клиент и ждет ввода команд
func Start() {
	cli := cmdInterface.New()
	conf := config.New()

	db, err := sql.Open("sqlite", "./sqlite-database.db")
	if err != nil {
		logrus.Panic(err)
	}
	defer db.Close()

	// TLS
	certConfig, err := getCert("./tls/ca.crt", "./tls/server.crt", "./tls/server.key")
	if err != nil {
		logrus.Panic(err)
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

	// CLI

	mc := NewMainController(cli, authUC, storeUC)
	mc.Register()

	cli.Run()
}

func getCert(ca string, crt string, key string) (*tls.Config, error) {
	caCert, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		logrus.Info("Not ok")
	}
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
