package qrcode

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteQrcodeStore struct {
	db *sql.DB
}

func NewSqliteQrcodeStore(path string) (*SqliteQrcodeStore, error) {
	var store SqliteQrcodeStore
	var err error

	// try opening the database
	store.db, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// create default table if it doesn't exists
	_, err = store.db.Exec("CREATE TABLE IF NOT EXISTS qrcodes (id INTEGER PRIMARY KEY, password VARCHAR(255), target VARCHAR(255))")
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *SqliteQrcodeStore) NewQrcode(ctx context.Context, password string, target string) (err error) {
	_, err = s.db.ExecContext(ctx, `
	INSERT INTO qrcodes (target, password) VALUES (?,?);
	`, target, password)

	return
}

func (s *SqliteQrcodeStore) GetQrcode(ctx context.Context, id int) (Qrcode, error) {
	var code Qrcode

	err := s.db.QueryRowContext(ctx, `
	SELECT  
		id, 
	    target,
	    password
	FROM qrcodes WHERE id=?;
	`, id).Scan(&code.Id, code.Target, code.Password)

	return code, err
}

func (s *SqliteQrcodeStore) GetQrcodeByPassword(ctx context.Context, password string) (Qrcode, error) {
	var code Qrcode

	err := s.db.QueryRowContext(ctx, `
	SELECT  
		id, 
	    target,
	    password
	FROM qrcodes WHERE password=?;
	`, password).Scan(&code.Id, code.Target, code.Password)

	return code, err
}

func (s *SqliteQrcodeStore) UpdateTargetUrl(ctx context.Context, id int, newTargetUrl string) error {

	_, err := s.db.ExecContext(ctx, `
	UPDATE qrcodes
	SET target = ?
	WHERE
		id = ?;
	`, newTargetUrl, id)

	return err
}
