package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	authHttpController "ya-gophkeeper-server/internal/auth/delivery/http"
	authPsqlRepository "ya-gophkeeper-server/internal/auth/repository/psql"
	authUseCase "ya-gophkeeper-server/internal/auth/usecase"
	"ya-gophkeeper-server/internal/config"
	middleware "ya-gophkeeper-server/internal/middlewares"

	storageHttpController "ya-gophkeeper-server/internal/storage/delivery/http"
	storagePsqlRepository "ya-gophkeeper-server/internal/storage/repository/psql"
	storageUseCase "ya-gophkeeper-server/internal/storage/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func Start() {
	router := chi.NewRouter()
	conf := config.New()

	router.Use(middleware.LogMiddleware)

	database, err := sql.Open("pgx", conf.DBConnect)
	if err != nil {
		logrus.Errorf("pgx init => %v", err)
		return
	}

	authRepo := authPsqlRepository.NewPsqlStorage(database)
	authUC := authUseCase.New(authRepo)
	authHttpController.NewHandler(conf, authUC).Register(router)

	storageRepo := storagePsqlRepository.NewPsqlStorage(database)
	storageUC := storageUseCase.New(storageRepo)
	storageHttpController.NewHandler(conf, storageUC).Register(router)

	logrus.Infof("Server started %s", conf.AddressHttp)

	caCert, err := os.ReadFile("./tls/ca.crt")
	if err != nil {
		logrus.Error(err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		logrus.Info("Not ok")
	}

	server := http.Server{
		Addr:    conf.AddressHttp,
		Handler: router,
		TLSConfig: &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	err = server.ListenAndServeTLS("./tls/server.crt", "./tls/server.key")
	// err = http.ListenAndServeTLS(conf.AddressHttp, "./ssl/server.crt", "./ssl/server.key", router)
	// err = http.ListenAndServe(conf.AddressHttp, router)
	if err != nil {
		logrus.Error(err)
	}
}
