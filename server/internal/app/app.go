package app

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	authHttpController "ya-gophkeeper-server/internal/auth/delivery/http"
	authPsqlRepository "ya-gophkeeper-server/internal/auth/repository/psql"
	authUseCase "ya-gophkeeper-server/internal/auth/usecase"
	"ya-gophkeeper-server/internal/config"

	storageHttpController "ya-gophkeeper-server/internal/storage/delivery/http"
	storagePsqlRepository "ya-gophkeeper-server/internal/storage/repository/psql"
	storageUseCase "ya-gophkeeper-server/internal/storage/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func Start() {
	router := chi.NewRouter()
	conf := config.New()

	database, err := sql.Open("pgx", conf.DBConnect)
	if err != nil {
		logrus.Errorf("pgx init => %v", err)
		return
	}

	authRepo := authPsqlRepository.NewPsqlStorage(database)
	authUC := authUseCase.New(authRepo)
	authHttpController.NewHandler(authUC).Register(router)

	storageRepo := storagePsqlRepository.NewPsqlStorage(database)
	storageUC := storageUseCase.New(storageRepo)
	storageHttpController.NewHandler(storageUC).Register(router)

	logrus.Infof("Server started %s", conf.AddressHttp)
	http.ListenAndServe(conf.AddressHttp, router)
}
