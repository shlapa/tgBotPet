package telegram

import (
	"log"
	"net/url"
	"strings"
)

const (
	Rnd   = "/rnd"
	Help  = "/help"
	Start = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("DO_COMMAND(%v, %v)", text, username)

	switch text {
	case Help:
		break
	case Rnd:
		break
	case Start:
		break
	default:
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Scheme != "" && u.Host != ""
}
