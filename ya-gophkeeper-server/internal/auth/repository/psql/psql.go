// Package psql пакет для работы с базой данных
package psql

import (
	"database/sql"

	"yandex-gophkeeper-server/internal/auth/entity"

	"github.com/sirupsen/logrus"
)

type psqlstorage struct {
	db *sql.DB
}

// NewPsqlStorage получить экземпляр структуры
func NewPsqlStorage(db *sql.DB) *psqlstorage {
	return &psqlstorage{
		db: db,
	}
}

// GetUser получить пользователя
func (repo *psqlstorage) GetUser(email string) (*entity.User, error) {
	user := entity.User{}
	row := repo.db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email)
	if err := row.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &user, nil
}

// CreateUser создать пользователя
func (repo *psqlstorage) CreateUser(email string, hashPassword string) (*entity.User, error) {
	row := repo.db.QueryRow("INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id", email, hashPassword)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	user := entity.User{
		ID:       id,
		Email:    email,
		Password: hashPassword,
	}
	return &user, nil
}
