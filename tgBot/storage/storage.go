package storage

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	hasher := md5.New()

	_, err := io.WriteString(hasher, p.URL)
	if err != nil {
		return "cant calculate hash URl: ", err
	}

	_, err = io.WriteString(hasher, p.UserName)
	if err != nil {
		return "cant calculate hash UserName: ", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
