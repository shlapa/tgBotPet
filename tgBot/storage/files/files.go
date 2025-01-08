package files

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
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

func (s Storage) LastLink(ctx context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't get last link", err) }()
	query := `SELECT * FROM tg_users WHERE "user" = $1 ORDER BY id_link DESC LIMIT 1`
	var linkDB string
	err = s.db.QueryRow(query, userName).Scan(&linkDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsLib.ErrNoSavedPage
		}
		return nil, err
	}

	page = &storage.Page{
		UserName: userName,
		URL:      linkDB,
	}
	return page, err
}

func (s Storage) SearchLink(ctx context.Context, p *storage.Page) (page *storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't pick link", err) }()
	query := `SELECT "link" FROM tg_users where "associations" = $1 ORDER BY RANDOM() DESC LIMIT 1`
	var linkDB string
	err = s.db.QueryRow(query, p.Associations).Scan(&linkDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsLib.ErrNoSavedPage
		}
		return nil, err
	}

	page = &storage.Page{
		URL: linkDB,
	}
	return page, err
}

func (s Storage) GetHistory(ctx context.Context, userName string) (pages []*storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't get history", err) }() // Обертка ошибки для улучшения контекста

	query := `SELECT "link" FROM tg_users WHERE "user" = $1`
	rows, err := s.db.QueryContext(ctx, query, userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pageList []*storage.Page

	for rows.Next() {
		var linkDB string
		if err := rows.Scan(&linkDB); err != nil {
			return nil, err
		}
		pageList = append(pageList, &storage.Page{
			UserName: userName,
			URL:      linkDB,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(pageList) == 0 {
		return nil, errorsLib.ErrNoSavedPage
	}
	return pageList, nil
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

func (s Storage) SaveAssociations(ctx context.Context, page *storage.Page) (err error) {
	defer func() { err = errorsLib.Wrap("can't save page", err) }()
	query := `UPDATE tg_users SET "associations" = $1 WHERE "user" = $2 AND "link" = $3`
	_, err = s.db.Exec(query, page.Associations, page.UserName, page.URL)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsLib.ErrNoSavedPage
		}
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

func (s Storage) IsLimit(ctx context.Context, p *storage.Page) (bool, error) {
	query := `SELECT COUNT(*) FROM tg_users where "user" = $1`
	var countDB string
	err := s.db.QueryRow(query, p.UserName).Scan(&countDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, nil
		}
	}
	number, err := strconv.Atoi(countDB)
	if err != nil {
		return false, err
	}
	if number <= 10 {
		return false, nil
	}
	return true, nil
}
