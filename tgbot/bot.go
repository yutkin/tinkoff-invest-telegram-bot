package tgbot

import (
	"log"
	"tinkoff-invest-telegram-bot/tinkoff"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TinkoffInvestmentsBot struct {
	TelegramgApi *tgbotapi.BotAPI
	TinkoffApi   *tinkoff.Api
	OwnerId      int
	WebHookToken string
}

func (bot *TinkoffInvestmentsBot) HandleCommandMessage(update *tgbotapi.Update) {
	if update.Message.From.ID != bot.OwnerId {
		log.Println("Unauthorized used_id", update.Message.From.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {

	case "portfolio":
		portfolio, err := bot.TinkoffApi.GetPortfolio()
		if err != nil {
			log.Println("Fail to get Portfolio", err)
			return
		}

		log.Println("Trying to get ticker names...")
		bot.TinkoffApi.FillPortfolioPositionsNames(&portfolio)

		log.Println("Trying to prettify portfolio...")
		prettyPortfolio, err := portfolio.Prettify(bot.TinkoffApi.PortfolioTemplate, bot.TinkoffApi.CurrencyConverter)
		log.Println("Successful prettifying!")

		if err != nil {
			log.Printf("Fail to prettify portfolio: %v\n", err)
			return
		}

		msg.Text = prettyPortfolio
		msg.ParseMode = tgbotapi.ModeHTML
	}

	_, err := bot.TelegramgApi.Send(msg)

	if err != nil {
		log.Printf("Fail to send message: %v", err)
	}
}
