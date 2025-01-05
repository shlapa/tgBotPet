package telegram

const msgHelp = `
Я твой верный оруженосец буду стараться сохранять все твои ссылки,
что ты возложишь на меня!

In order to save the page, just send me al link to it.

In order to get a random page from your list, send me command /rnd.
Caution! After that, this page will be removed from your list!`

const msgHello = "Приветствую тебя воин! ⚔️\n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages"
	msgSaved          = "Saved"
	msgAlreadyExists  = "You have already have this page in your list"
	msgHaveNotLinked  = "You have not linked to this page"
)
