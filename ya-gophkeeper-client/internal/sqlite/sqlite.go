package sqlite

import (
	"database/sql"

	"ya-gophkeeper-client/internal/entity"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLite {
	return &SQLite{
		db: db,
	}
}

func (sq *SQLite) Migration() error {
	_, err := sq.db.Exec(`
		CREATE TABLE IF NOT EXISTS data(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB NOT NULL,
			meta VARCHAR(1024)
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func (sq *SQLite) Close() error {
	return sq.db.Close()
}

func (sq *SQLite) Save(item *entity.ItemDto) error {
	_, err := sq.db.Exec("INSERT INTO data (data, meta) VALUES (?, ?)", item.Data, item.Meta)
	return err
}

func (sq *SQLite) GetAll() (*[]entity.ItemDto, error) {
	rows, err := sq.db.Query("SELECT data, meta FROM data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]entity.ItemDto, 0)
	for rows.Next() {
		item := entity.ItemDto{}
		if err := rows.Scan(&item.Data, &item.Meta); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &items, err
}
