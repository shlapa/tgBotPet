package telegram

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"strings"
	"tgBot/lib/errorsLib"
	"tgBot/storage"
)

const (
	Rnd        = "/rnd"
	Help       = "/help"
	Start      = "/start"
	Delete     = "/delete"
	LastLink   = "/get_last_link"
	SearchLink = "/search_link"
	GetHistory = "/get_history"
	DeleteALl  = "/delete_all"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	pageLastLink := &storage.Page{
		UserName: username,
	}
	log.Printf("DO_COMMAND(%v, %v)", text, username)

	//Возможность ассуждать l;)
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	ban1 := os.Getenv("BAN_B1")
	ban2 := os.Getenv("BAN_B2")
	ban3 := os.Getenv("BAN_B3")
	words := []string{ban1, ban2, ban3}

	if isAddCmd(text) {
		for _, word := range words {
			if strings.Contains(text, word) {
				if err := p.tg.SendMessage(chatID, "Богохульство! Но я сохраню !😤"); err != nil {
					return err
				}
			}
		}
		return p.savePage(text, chatID, username)
	} else if isAddText(text) {
		for _, word := range words {
			if strings.Contains(text, word) {
				if err := p.tg.SendMessage(chatID, "Богохульство! Но я сохраню !😤"); err != nil {
					return err
				}
			}
		}
		return p.processAssociations(chatID, text)
	}

	err := p.clearLastLink(chatID)
	if err != nil {
		return err
	}

	if strings.HasPrefix(text, Delete) && len(strings.TrimPrefix(text, Delete)) > 0 && strings.TrimPrefix(text, Delete)[0] == ' ' {
		space := strings.TrimSpace(strings.TrimPrefix(text, Delete))
		if space == "" {
			return p.tg.SendMessage(chatID, "Пожалуйста, укажи ссылку, что желаешь уничтожить. 🔗💀")
		}
		pageLastLink.URL = space
		return p.Remove(chatID, pageLastLink)
	}

	if strings.HasPrefix(text, SearchLink) {
		space := strings.TrimSpace(strings.TrimPrefix(text, SearchLink))
		if space == "" {
			return p.tg.SendMessage(chatID, "Милорд, мне нужен след, чтобы начать поиски. 🔍 Укажите его в формате: /search_link след1, след2 и так далее. 🗺️")
		}
		pageLastLink.Associations = space
		return p.searchLink(chatID, pageLastLink)
	}

	switch text {
	case Help:
		return p.sendHelp(chatID)
	case Rnd:
		return p.sendRandom(chatID, username)
	case Start:
		return p.sendHello(chatID, username)
	case LastLink:
		return p.getLastLink(chatID, username)
	case GetHistory:
		return p.getHistory(chatID, username)
	case DeleteALl:
		return p.deleteAll(chatID, username)

	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) deleteAll(chatID int, username string) (err error) {
	defer func() { err = errorsLib.Wrap("cantDeleteAll", err) }()
	err = p.storage.RemoveAll(context.Background(), username)
	if err != nil {
		return p.tg.SendMessage(chatID, "Эта ссылка исчезла в туманном мире... 👻")
	}
	return p.tg.SendMessage(chatID, "Похоже, твои свитки исчезли в бездне времени... ⏳")
}

func (p *Processor) getLastLink(chatID int, username string) (err error) {
	defer func() { err = errorsLib.Wrap("cantGetLastLink", err) }()

	page, err := p.storage.LastLink(context.Background(), username)
	if err != nil {
		if errors.Is(err, errorsLib.ErrNoSavedPage) {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		} else {
			return err
		}
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	p.lastLink[chatID] = page

	return nil
}

func (p *Processor) searchLink(chatID int, pageLastLink *storage.Page) (err error) {
	defer func() { err = errorsLib.Wrap("cantGetLastLink", err) }()
	page, err := p.storage.SearchLink(context.Background(), pageLastLink)
	if err != nil {
		if errors.Is(err, errorsLib.ErrNoSavedPage) {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		} else {
			return err
		}
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	p.lastLink[chatID] = page
	return nil
}

func (p *Processor) getHistory(chatID int, username string) error {
	pages, err := p.storage.GetHistory(context.Background(), username)
	if err != nil {
		if errors.Is(err, errorsLib.ErrNoSavedPage) {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		} else {
			return err
		}
	}

	if len(pages) == 0 {
		return p.tg.SendMessage(chatID, msgHaveNotLinked)
	}

	for _, page := range pages {
		if err := p.tg.SendMessage(chatID, page.URL); err != nil {
			return err
		}
	}

	p.lastLink[chatID] = pages[len(pages)-1]

	return nil
}

func (p *Processor) Remove(chatID int, pageLastLink *storage.Page) error {
	err := p.storage.Remove(context.Background(), pageLastLink)
	if err != nil {
		return p.tg.SendMessage(chatID, "Эта ссылка исчезла в туманном мире... 👻")
	}
	return p.tg.SendMessage(chatID, "Ты избавился от свитка, как рыцарь от старого оружия. ⚔️")
}

func (p *Processor) savePage(textURL string, chatID int, username string) (err error) {
	defer func() { err = errorsLib.Wrap("cantSavePage", err) }()

	page := &storage.Page{
		URL:          textURL,
		UserName:     username,
		Associations: "",
	}

	isLimit, err := p.storage.IsLimit(context.Background(), page)
	if err != nil {
		return err
	}
	if isLimit {
		return p.tg.SendMessage(chatID, msgLimitExceeded)
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

	if err = p.tg.SendMessage(chatID, "Желаешь добавить ассоциации к этому свитку? ✍️"); err != nil {
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
		return p.tg.SendMessage(chatID, "Рад бы поболтать, о славный рыцарь, но королевство в опасности, а работы невпроворот! 🏰⚔️ Попробуй задать команду, чтобы помочь делу!\n")
	}

	page.Associations = input

	if err := p.storage.SaveAssociations(context.Background(), page); err != nil {
		return errorsLib.Wrap("cantSaveAssociations", err)
	}

	return p.tg.SendMessage(chatID, "Ассоциации успешно добавлены. ✨")
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = errorsLib.Wrap("cantPickPage", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil {
		if errors.Is(err, errorsLib.ErrNoSavedPage) {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		} else {
			return err
		}
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
	return true
}

func (p *Processor) clearLastLink(chatID int) error {
	delete(p.lastLink, chatID)
	return p.tg.SendMessage(chatID, "Последняя ссылка была удалена. Больше я не помню её. 😶‍🌫️")
}
