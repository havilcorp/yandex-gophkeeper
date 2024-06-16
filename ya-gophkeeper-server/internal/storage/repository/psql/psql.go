package psql

import (
	"database/sql"

	"yandex-gophkeeper-server/internal/storage/entity"
)

type psqlstorage struct {
	db *sql.DB
}

func NewPsqlStorage(db *sql.DB) *psqlstorage {
	return &psqlstorage{
		db: db,
	}
}

// TODO: AddContext
func (repo *psqlstorage) Save(userID int, dto *entity.CreateDto) error {
	_, err := repo.db.Exec("INSERT INTO data (user_id, data, meta) VALUES ($1, $2, $3)", userID, dto.Data, dto.Meta)
	return err
}

func (repo *psqlstorage) GetById(id int) (*entity.Item, error) {
	row := repo.db.QueryRow("SELECT id, user_id, data, meta FROM data WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var item entity.Item
	if err := row.Scan(&item.ID, &item.UserId, &item.Data, &item.Meta); err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *psqlstorage) GetAll(userID int) (*[]entity.Item, error) {
	rows, err := repo.db.Query("SELECT id, user_id, data, meta FROM data WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]entity.Item, 0)
	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(&item.ID, &item.UserId, &item.Data, &item.Meta); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &items, nil
}

func (repo *psqlstorage) Remove(id int) error {
	_, err := repo.db.Exec("DELETE FROM data WHERE id = $1", id)
	return err
}
