# Project Name

Это проект, который помогает управлять ссылками и сохранять их для последующего использования. Бот поддерживает различные команды для сохранения, поиска и удаления ссылок, а также добавления ассоциаций. Изначально был разработан для работы с локальным сохранением данных, но впоследствии был расширен для работы с базой данных.

## Структура веток

- **main**: Эта ветка содержит первоначальную версию бота, которая работает с локальным сохранением ссылок.
- **bdaccess**: В этой ветке был расширен функционал, и теперь бот работает с базой данных для хранения ссылок и ассоциаций.

## Функционал
-    /start
    Запуск бота. Инициализация взаимодействия с пользователем и начало работы с функционалом.
-   /help
    Получение справки по доступным командам. Описание возможностей бота для удобства пользователя.
-  /delete [ссылка]
    Удаление сохраненной ссылки из базы данных. Пользователь указывает ссылку, которую необходимо удалить.
-  /get_last_link
    Получение последней сохраненной пользователем ссылки. Бот возвращает ссылку, добавленную в базу данных последней.
-  /search_link [ключи поиска]
    Поиск сохраненных ссылок по ключевым словам. Пользователь указывает ключи, и бот возвращает соответствующие результаты.
-  /get_history
    Получение истории всех сохраненных ссылок. Бот выводит список всех ссылок, добавленных пользователем.
-  /search_traces
    Поиск ссылок по опсианию.
-  /delete_all
    Полная очистка базы данных. Удаление всех сохраненных пользователем ссылок.
