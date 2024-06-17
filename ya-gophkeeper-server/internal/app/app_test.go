package app

import (
	"testing"
	"time"

	"yandex-gophkeeper-server/internal/config"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func Test_startServer(t *testing.T) {
	conf := config.New()
	conf.CACrt = "../../tls/ca.crt"
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	serverHttp, grpcListener, err := startServer(conf, db)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second * 2)

	if err := serverHttp.Close(); err != nil {
		logrus.Error(err)
	}

	if err := grpcListener.Close(); err != nil {
		logrus.Error(err)
	}
}
