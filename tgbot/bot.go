package tgbot

import (
	"fmt"
	"log"
	"tinkoff-investments-telegram-bot/tinkoff"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TinkoffInvestmentsBot struct {
	TelegramgApi *tgbotapi.BotAPI
	TinkoffApi   *tinkoff.Api
	OwnerId int
}

func (bot *TinkoffInvestmentsBot) HandleCommandMessage(update *tgbotapi.Update) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {

	case "portfolio":
		portfolio, err := bot.TinkoffApi.GetPortfolio()
		if err == nil {
			msg.Text = portfolio.Prettify()
			msg.ParseMode = tgbotapi.ModeHTML
		} else {
			msg.Text = fmt.Sprintf("%v", err)
		}
	}

	_, err := bot.TelegramgApi.Send(msg)

	if err != nil {
		log.Printf("Fail to send message: %v", err)
	}
}
