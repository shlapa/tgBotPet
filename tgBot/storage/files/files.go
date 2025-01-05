package files

import (
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"tgBot/lib/errorsLib"
	"tgBot/storage"
	"time"
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
	//err = s.db.QueryRow(query, page.UserName, page.URL, page.Associations).Scan()
	_, err = s.db.Exec(query, page.UserName, page.URL, page.Associations)
	if err != nil {
		return err
	}
	return nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) PickRandom(ctx context.Context, userName string) (page *storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't pick random page: ", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		// Проверяем, существует ли директория
		if os.IsNotExist(err) {
			// Если директории нет, создаем её
			if mkErr := os.MkdirAll(path, os.ModePerm); mkErr != nil {
				return nil, fmt.Errorf("failed to create directory %s: %w", path, mkErr)
			}
			// Возвращаем ошибку, что директория создана, но файлов нет
			return nil, errors.New("no files found, created empty directory")
		}
		// Если ошибка не связана с отсутствием директории, возвращаем её
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	if len(files) == 0 {
		return nil, errorsLib.ErrNoSavedPage
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(files))

	file := files[r]
	if file == nil {
		return nil, errors.New("selected file is nil")
	}

	path = filepath.Join(path, file.Name())

	return s.DecodePage(ctx, filepath.Join(path))
}

func (s Storage) Remove(ctx context.Context, page *storage.Page) (err error) {
	if page == nil {
		return errorsLib.Wrap("page is nil", errors.New("page is nil"))
	}
	fileName, err := fileName(page)
	if err != nil {
		return errorsLib.Wrap("can't delete random page: ", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fileName)
	if err = os.Remove(path); err != nil {
		return errorsLib.Wrap("can't remove file: ", err)
	}
	path = filepath.Join(s.basePath, page.UserName, path)

	if err = os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file: %s", path)
		return errorsLib.Wrap(msg, err)
	}

	return nil
}

func (s Storage) DecodePage(ctx context.Context, filePath string) (page *storage.Page, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errorsLib.Wrap("can't open file: ", err)
	}
	defer func() { _ = file.Close() }()

	var p storage.Page

	if err = gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, errorsLib.Wrap("can't decode file: ", err)
	}
	return &p, nil
}

func (s Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, errorsLib.Wrap("can't generate file name: ", err)
	}

	// Собираем полный путь к файлу
	path := filepath.Join(s.basePath, p.UserName, fileName)
	log.Printf("Checking existence of file at path: %s\n", path)

	// Проверяем, существует ли директория
	dirPath := filepath.Join(s.basePath, p.UserName)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return false, errorsLib.Wrap(fmt.Sprintf("can't create directory: %s", dirPath), err)
	}

	// Проверяем, существует ли сам файл
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("File does not exist:", path)
		return false, nil // Файл не найден, это не ошибка
	} else if err != nil {
		return false, errorsLib.Wrap(fmt.Sprintf("error checking file existence at path: %s", path), err)
	}

	return true, nil // Файл существует
}
