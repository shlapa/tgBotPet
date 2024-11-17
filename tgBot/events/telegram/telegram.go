package telegram

import (
	"tgBot/clients/telegram"
	"tgBot/events"
	"tgBot/lib/errorsLib"
	"tgBot/storage"
	"tgBot/storage/files"
)

type Processor struct {
	tg       *telegram.Client
	offset   int
	storage  storage.Storage
	lastLink map[int]*storage.Page
}

type Meta struct {
	chatID   int
	userName string
}

func New(client *telegram.Client, storageFiles files.Storage) *Processor {
	return &Processor{
		tg:       client,
		storage:  storageFiles,
		lastLink: make(map[int]*storage.Page),
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errorsLib.Wrap("can't fetch updates", err)
	}

	if len(update) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(update))

	for _, up := range update {
		result = append(result, event(up))
	}

	p.offset = update[len(result)-1].UpdateId + 1

	return result, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errorsLib.Wrap("can't process event", errorsLib.ErrUnknownProcess)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errorsLib.Wrap("can't process event", err)
	}

	if err := p.doCmd(event.Text, meta.chatID, meta.userName); err != nil {
		return errorsLib.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errorsLib.Wrap("can't get  Meta", errorsLib.ErrorTypeMeta)
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)

	res := events.Event{
		Type: uType,
		Text: fetchText(u),
	}

	if uType == events.Message {
		res.Meta = Meta{
			chatID:   u.Message.Chat.Id,
			userName: u.Message.From.Username,
		}
	}
	return res
}

func fetchText(u telegram.Update) string {
	if u.Message != nil {
		return u.Message.Text
	}
	return ""
}

func fetchType(u telegram.Update) events.Type {
	if u.Message != nil {
		return events.Message
	}
	return events.Unknown
}
