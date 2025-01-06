package telegram

const msgHelp = `
Я твой верный оруженосец, клятвенно обручаю свою службу в сохранении всех твоих важных свитков 📜 и ссылок 🔗, что ты возложишь на мои плечи!

Чтобы сохранить страницу, просто пришли мне её ссылку, и я сохраню её для тебя в нашей великой книге 📚.

Для того чтобы извлечь одну случайную страницу из твоего списка, возьми в руку меч ⚔️ и изречь команду /rnd.
Будь осторожен! После этого свиток будет удалён из списка, словно исчезнет в небытие! 💨
`

const msgHello = "Здравствуй, благородный воин! ⚔️\n\n" + msgHelp

const (
	msgUnknownCommand = "Твоя команда мне неведома, милорд. 🤔🔮"
	msgNoSavedPages   = "Ты не связал ни одного свитка с нашим кодексом, мой рыцарь. 📜❌"
	msgSaved          = "Твой свиток был сохранён в великой книге. 📚✨"
	msgAlreadyExists  = "Ты уже оставил эту страницу в своём списке. Повторно не возлагаешь! 🔁⚔️"
	msgHaveNotLinked  = "Ты не связал эту страницу с нами. Не опозорь себя! 😓🔗"
)
