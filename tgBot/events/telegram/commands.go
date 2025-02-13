package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"tgBot/lib/errorsLib"
	"tgBot/storage"
)

const (
	Rnd           = "/rnd"
	Help          = "/help"
	Start         = "/start"
	Delete        = "/delete"
	History       = "/history"
	DeleteHistory = "/hd"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("DO_COMMAND(%v, %v)", text, username)

	if isAddCmd(text) {
		return p.savePage(text, chatID, username)
	}

	switch text {
	case Help:
		return p.sendHelp(chatID)
	case Rnd:
		return p.sendRandom(chatID, username)
	case Start:
		return p.sendHello(chatID, username)
	case Delete:
		page := p.lastLink[chatID]
		return p.storage.Remove(context.Background(), page)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(textURL string, chatID int, username string) (err error) {
	defer func() { err = errorsLib.Wrap("cantSavePage", err) }()

	page := &storage.Page{
		URL:      textURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err = p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err = p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	p.lastLink[chatID] = page

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = errorsLib.Wrap("cantSavePage", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && errors.Is(err, errorsLib.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgHaveNotLinked)
	}
	if errors.Is(err, errorsLib.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	p.lastLink[chatID] = page

	return nil
}

func (p *Processor) sendHelp(chatID int) (err error) {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int, username string) (err error) {
	return p.tg.SendMessage(chatID, username+"! "+msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Scheme != "" && u.Host != ""
}
