package files

import (
	"encoding/gob"
	"errors"
	"fmt"
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
}

func NewStorage(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() {
		if err != nil {
			err = errorsLib.Wrap("can't save: ", err)
		}
	}()

	// Определяем путь к файлу
	filePath := filepath.Join(s.basePath, page.UserName)

	// Создаем необходимые директории
	if err = os.MkdirAll(filepath.Dir(filePath), defParam); err != nil {
		return err
	}

	// Генерируем имя файла
	fileName, err := fileName(page)
	if err != nil {
		return err
	}

	// Полный путь к файлу
	filePath = filepath.Join(filePath, fileName)

	// Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	err = gob.NewEncoder(file).Encode(page)
	if err != nil {
		return err
	}

	return nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = errorsLib.Wrap("can't pick random page: ", err) }()
	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errorsLib.ErrNoSavedPage
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(files))

	file := files[r]

	return s.DecodePage(filepath.Join(file.Name()))
}

func (s Storage) Remove(page *storage.Page) (err error) {
	filename := filepath.Join(s.basePath, page.UserName)
	if err = os.Remove(filename); err != nil {
		return errorsLib.Wrap("can't remove file: ", err)
	}
	path := filepath.Join(s.basePath, page.UserName, filename)

	if err = os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file: %s", path)
		return errorsLib.Wrap(msg, err)
	}

	return nil
}

func (s Storage) DecodePage(filePath string) (page *storage.Page, err error) {
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

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, errorsLib.Wrap("can't check if file exists: ", err)
	}
	path := filepath.Join(s.basePath, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, err
	case err != nil:
		msg := fmt.Sprintf("can't check if file exists: %s", path)
		return false, errorsLib.Wrap(msg, err)
	}

	return true, nil
}
