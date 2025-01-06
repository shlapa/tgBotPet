package files

import (
	"context"
	"database/sql"
	"errors"
	"tgBot/lib/errorsLib"
	"tgBot/storage"
)

const (
	defParam = 0774
)

type Storage struct {
	basePath string
	db       *sql.DB
}

func NewStorage(basePath string, db *sql.DB) Storage {
	return Storage{
		basePath: basePath,
		db:       db,
	}
}

func (s Storage) Save(ctx context.Context, page *storage.Page) (err error) {
	defer func() { err = errorsLib.Wrap("can't save page", err) }()
	query := `INSERT INTO tg_users ("user", "link", "associations") VALUES ($1, $2, $3)`
	_, err = s.db.Exec(query, page.UserName, page.URL, page.Associations)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(ctx context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't pick random page: ", err) }()
	query := `SELECT "link" FROM tg_users where "user" = $1 ORDER BY RANDOM() DESC LIMIT 1`
	var linkDB string
	err = s.db.QueryRow(query, userName).Scan(&linkDB)
	if err != nil {
		return nil, err
	}

	page = &storage.Page{
		UserName: userName,
		URL:      linkDB,
	}
	return page, err
}

func (s Storage) Remove(ctx context.Context, page *storage.Page) (err error) {
	defer func() { err = errorsLib.Wrap("can't save page", err) }()
	query := `DELETE FROM tg_users WHERE "user" = $1 AND "link" = $2`
	_, err = s.db.Exec(query, page.UserName, page.URL)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	query := `SELECT "user", "link" FROM tg_users where "user" = $1 and link = $2`
	var userDB, linkDB string
	err := s.db.QueryRow(query, p.UserName, p.URL).Scan(&userDB, &linkDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, nil
		}
	}
	return true, nil
}
