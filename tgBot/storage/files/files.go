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
	defer func() { err = errorsLib.Wrap("can't get last link for user "+userName, err) }()

	query := `SELECT "link" FROM tg_users WHERE "user" = $1 AND "user" = $2 ORDER BY id_link DESC LIMIT 1`
	var linkDB string

	err = s.db.QueryRow(query, userName, userName).Scan(&linkDB)
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
	return page, nil
}

func (s Storage) SearchLink(ctx context.Context, p *storage.Page) (page *storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't pick link for associations "+p.Associations, err) }()
	query := `SELECT "link" FROM tg_users WHERE "associations" = $1 AND "user" = $2 ORDER BY RANDOM() DESC LIMIT 1`
	var linkDB string
	err = s.db.QueryRow(query, p.Associations, p.UserName).Scan(&linkDB)
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
	defer func() { err = errorsLib.Wrap("can't get history for user "+userName, err) }()
	query := `SELECT "link" FROM tg_users WHERE "user" = $1 AND "user" = $2`
	rows, err := s.db.QueryContext(ctx, query, userName, userName)
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
	defer func() { err = errorsLib.Wrap("can't save page for user "+page.UserName, err) }()
	query := `INSERT INTO tg_users ("user", "link", "associations") VALUES ($1, $2, $3)`
	_, err = s.db.Exec(query, page.UserName, page.URL, page.Associations)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) SaveAssociations(ctx context.Context, page *storage.Page) (err error) {
	defer func() { err = errorsLib.Wrap("can't save associations for user "+page.UserName, err) }()
	query := `UPDATE tg_users SET "associations" = $1 WHERE "user" = $2 AND "link" = $3`
	_, err = s.db.Exec(query, page.Associations, page.UserName, page.URL)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(ctx context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't pick random page for user "+userName, err) }()
	query := `SELECT "link" FROM tg_users WHERE "user" = $1 AND "user" = $2 ORDER BY RANDOM() DESC LIMIT 1`
	var linkDB string
	err = s.db.QueryRow(query, userName, userName).Scan(&linkDB)
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
	defer func() { err = errorsLib.Wrap("can't remove page for user "+page.UserName, err) }()
	query := `DELETE FROM tg_users WHERE "user" = $1 AND "link" = $2 AND "user" = $3`
	_, err = s.db.Exec(query, page.UserName, page.URL, page.UserName)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) RemoveAll(ctx context.Context, userName string) (err error) {
	defer func() { err = errorsLib.Wrap("can't remove all pages for user "+userName, err) }()
	query := `DELETE FROM tg_users WHERE "user" = $1 AND "user" = $2`
	_, err = s.db.Exec(query, userName, userName)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	query := `SELECT "user", "link" FROM tg_users WHERE "user" = $1 AND "link" = $2 AND "user" = $3`
	var userDB, linkDB string
	err := s.db.QueryRow(query, p.UserName, p.URL, p.UserName).Scan(&userDB, &linkDB)
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
	query := `SELECT COUNT(*) FROM tg_users WHERE "user" = $1 AND "user" = $2`
	var countDB string
	err := s.db.QueryRow(query, p.UserName, p.UserName).Scan(&countDB)
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
