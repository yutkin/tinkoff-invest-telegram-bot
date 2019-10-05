module tinkoff-investments-telegram-bot

require (
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	tinkoff-investments-telegram-bot/tgbot v0.0.0
	tinkoff-investments-telegram-bot/tinkoff v0.0.0
)

replace tinkoff-investments-telegram-bot/tinkoff => ./tinkoff

replace tinkoff-investments-telegram-bot/tgbot => ./tgbot

go 1.11
