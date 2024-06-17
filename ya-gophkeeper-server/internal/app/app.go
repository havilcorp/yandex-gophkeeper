package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	authHttpController "yandex-gophkeeper-server/internal/auth/delivery/http"
	authPsqlRepository "yandex-gophkeeper-server/internal/auth/repository/psql"
	authUseCase "yandex-gophkeeper-server/internal/auth/usecase"
	"yandex-gophkeeper-server/internal/config"
	middleware "yandex-gophkeeper-server/internal/middlewares"

	storageGRPCController "yandex-gophkeeper-server/internal/storage/delivery/grpc"
	storageHttpController "yandex-gophkeeper-server/internal/storage/delivery/http"
	storagePsqlRepository "yandex-gophkeeper-server/internal/storage/repository/psql"
	storageUseCase "yandex-gophkeeper-server/internal/storage/usecase"

	pb "github.com/havilcorp/yandex-gophkeeper-proto/save"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func startServer(conf *config.Config, database *sql.DB) (*http.Server, net.Listener, error) {
	router := chi.NewRouter()

	router.Use(middleware.LogMiddleware)

	authRepo := authPsqlRepository.NewPsqlStorage(database)
	authUC := authUseCase.New(authRepo)
	authHttpController.NewHandler(conf, authUC).Register(router)

	storageRepo := storagePsqlRepository.NewPsqlStorage(database)
	storageUC := storageUseCase.New(storageRepo)
	storageHttpController.NewHandler(conf, storageUC).Register(router)

	caCert, err := os.ReadFile(conf.CACrt)
	if err != nil {
		return nil, nil, fmt.Errorf("os.ReadFile: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, nil, fmt.Errorf("caCertPool.AppendCertsFromPEM: %w", err)
	}

	serverHttp := http.Server{
		Addr:    conf.AddressHttp,
		Handler: router,
		TLSConfig: &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	grpcListener, err := net.Listen("tcp", conf.AddressGRPC)
	if err != nil {
		return nil, nil, fmt.Errorf("net.Listen: %w", err)
	}
	cred, err := credentials.NewServerTLSFromFile(conf.ServerCrt, conf.ServerKey)
	serverGRPC := grpc.NewServer(grpc.Creds(cred), grpc.ChainUnaryInterceptor(middleware.AuthGRPCMiddleware(conf.JWTKey)))
	pb.RegisterSaveServer(serverGRPC, storageGRPCController.NewHandler(conf, storageUC))

	go func() {
		logrus.Printf("Сервер gRPC начал работу по адресу: %s\n", conf.AddressGRPC)
		if err := serverGRPC.Serve(grpcListener); err != nil {
			logrus.Error(err)
		}
		logrus.Printf("Сервер gRPC прекратил работу")
	}()

	go func() {
		logrus.Printf("Сервер HTTP начал работу по адресу: %s\n", conf.AddressHttp)
		err := serverHttp.ListenAndServeTLS(conf.ServerCrt, conf.ServerKey)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Printf("Сервер HTTP прекратил работу")
	}()

	return &serverHttp, grpcListener, nil
}

func Start() error {
	conf := config.New()
	database, err := sql.Open("pgx", conf.DBConnect)
	if err != nil {
		logrus.Errorf("pgx init => %v", err)
		return err
	}
	defer database.Close()

	serverHttp, grpcListener, err := startServer(conf, database)
	if err != nil {
		serverHttp.Close()
		grpcListener.Close()
		return err
	}

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-terminateSignals

	if err := serverHttp.Close(); err != nil {
		logrus.Error(err)
	}

	if err := grpcListener.Close(); err != nil {
		logrus.Error(err)
	}

	time.Sleep(time.Second * 2)

	return nil
}
