package psql

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"yandex-gophkeeper-server/internal/storage/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_psqlstorage_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sql := "INSERT INTO data (user_id, data, meta) VALUES ($1, $2, $3)"
	mock.ExpectExec(regexp.QuoteMeta(sql)).WithArgs(1, []byte("data"), "meta").WillReturnResult(driver.ResultNoRows)

	repo := NewPsqlStorage(db)
	err = repo.Save(1, &entity.CreateDto{
		Data: []byte("data"),
		Meta: "meta",
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_psqlstorage_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sql := "SELECT id, user_id, data, meta FROM data WHERE user_id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "data", "meta"}).
			AddRow(1, 1, []byte("data"), "meta"))
	mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs(2).WillReturnError(errors.New(""))

	type args struct {
		userID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				userID: 1,
			},
			wantErr: false,
		},
		{
			name: "error query",
			args: args{
				userID: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewPsqlStorage(db)
			_, err := repo.GetAll(tt.args.userID)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
