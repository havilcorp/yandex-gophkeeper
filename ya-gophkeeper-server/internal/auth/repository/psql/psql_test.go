package psql

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_psqlstorage_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sql := "SELECT id, email, password FROM users WHERE email = $1"
	mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs("mail@mail.ru").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(1, "mail@mail.ru", "1a1dc91c907325c69271ddf0c944bc72"))
	mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs("err1@mail.ru").WillReturnError(errors.New(""))
	mock.ExpectQuery(regexp.QuoteMeta(sql)).WithArgs("err2@mail.ru").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			RowError(1, errors.New("")))

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				email: "mail@mail.ru",
			},
			wantErr: false,
		},
		{
			name: "error query",
			args: args{
				email: "err1@mail.ru",
			},
			wantErr: true,
		},
		{
			name: "error scan",
			args: args{
				email: "err2@mail.ru",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewPsqlStorage(db)
			_, err := repo.GetUser(tt.args.email)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func Test_psqlstorage_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sql := "INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id"
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs("mail@mail.ru", "1a1dc91c907325c69271ddf0c944bc72").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs("err1@mail.ru", "1a1dc91c907325c69271ddf0c944bc72").
		WillReturnError(errors.New(""))
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs("err2@mail.ru", "1a1dc91c907325c69271ddf0c944bc72").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("")))

	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				email:    "mail@mail.ru",
				password: "1a1dc91c907325c69271ddf0c944bc72",
			},
			wantErr: false,
		},
		{
			name: "error query",
			args: args{
				email:    "err1@mail.ru",
				password: "1a1dc91c907325c69271ddf0c944bc72",
			},
			wantErr: true,
		},
		{
			name: "err scan",
			args: args{
				email:    "err2@mail.ru",
				password: "1a1dc91c907325c69271ddf0c944bc72",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewPsqlStorage(db)
			_, err := repo.CreateUser(tt.args.email, tt.args.password)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
