package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
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
