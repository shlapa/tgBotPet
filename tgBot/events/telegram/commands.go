package telegram

import (
	"context"
	"errors"
	"fmt"
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
	LastLink      = "/get_last_link"
	DeleteHistory = "/deleteLink"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("DO_COMMAND(%v, %v)", text, username)

	// Если это ссылка
	if isAddCmd(text) {
		return p.savePage(text, chatID, username)
	} else if isAddText(text) {
		return p.processAssociations(chatID, text)
	}

	// Если это ассоциации

	switch text {
	case Help:
		return p.sendHelp(chatID)
	case Rnd:
		return p.sendRandom(chatID, username)
	case Start:
		return p.sendHello(chatID, username)
	case Delete:
		page := p.lastLink[chatID]
		return p.Remove(chatID, page)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) Remove(id int, page *storage.Page) error {

	return nil
}

func (p *Processor) savePage(textURL string, chatID int, username string) (err error) {
	defer func() { err = errorsLib.Wrap("cantSavePage", err) }()

	page := &storage.Page{
		URL:          textURL,
		UserName:     username,
		Associations: "",
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

	if err = p.tg.SendMessage(chatID, "Добавить ассоциации для сохраненной ссылки ?:"); err != nil {
		return err
	}

	return nil
}

func (p *Processor) AddAssociations(chatID int, input string) error {
	return p.processAssociations(chatID, input)
}

func (p *Processor) processAssociations(chatID int, input string) error {
	page, ok := p.lastLink[chatID]
	if !ok {
		return p.tg.SendMessage(chatID, "Не удалось найти последнюю ссылку. Попробуйте снова.")
	}

	// Сохранение ассоциаций
	page.Associations = input

	if err := p.storage.Save(context.Background(), page); err != nil {
		return errorsLib.Wrap("cantSaveAssociations", err)
	}

	return p.tg.SendMessage(chatID, "Ассоциации успешно сохранены.")
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

func isAddText(text string) bool {
	return isText(text)
}

func isText(text string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}
	if isURL(text) {
		return false
	}
	if strings.HasPrefix(text, "/") {
		return false
	}
	fmt.Printf("")
	return true
}
