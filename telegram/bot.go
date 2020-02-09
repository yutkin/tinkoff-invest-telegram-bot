package telegram

import (
	"fmt"
	"log"
	"time"
	"investbot/currency"
	"investbot/tinkoff"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	telegramApi       *tgbotapi.BotAPI
	tinkoffApi        *tinkoff.Api
	currencyConverter currency.Converter
	ownerId           int
	WebHookToken      string
}

var commandsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/broker"),
		tgbotapi.NewKeyboardButton("/iis"),
	),
)

func New(telegramBotToken, webHookToken, currencyConvertToken string, botOwnerId int, tinkoffApi *tinkoff.Api) (*Bot, error) {
	bot := &Bot{
		tinkoffApi:        tinkoffApi,
		ownerId:           botOwnerId,
		WebHookToken:      webHookToken,
		currencyConverter: currency.New(currencyConvertToken, 5*time.Second),
	}

	var err error
	bot.telegramApi, err = tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		return nil, err
	}

	log.Println("Authorized on account", bot.telegramApi.Self.UserName)

	return bot, nil
}

func (bot *Bot) sendPortfolioAsMessage(chatId int64, portfolio tinkoff.Portfolio) error {
	prettifiedPortfolio, err := portfolio.Prettify(bot.currencyConverter)

	if err != nil {
		return fmt.Errorf("fail to prettify portfolio: %v", err)
	}

	msg := tgbotapi.NewMessage(chatId, prettifiedPortfolio)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = commandsKeyboard

	_, err = bot.telegramApi.Send(msg)

	if err != nil {
		return fmt.Errorf("fail to send message: %v", err)
	}

	return nil
}

func (bot *Bot) HandleCommandMessage(update *tgbotapi.Update) {
	if update.Message.From.ID != bot.ownerId {
		log.Println("Unauthorized used_id", update.Message.From.ID)
		return
	}

	switch update.Message.Command() {
	case "broker":
		portfolio, err := bot.tinkoffApi.GetBrokerPortfolio()
		if err != nil {
			log.Println("Fail to fetch portfolio:", err)
			return
		}

		if err := bot.sendPortfolioAsMessage(update.Message.Chat.ID, portfolio); err != nil {
			log.Println("Fail to send message:", err)
		}

	case "iis":
		portfolio, err := bot.tinkoffApi.GetIisPortfolio()
		if err != nil {
			log.Println("Fail to fetch portfolio:", err)
			return
		}

		if err := bot.sendPortfolioAsMessage(update.Message.Chat.ID, portfolio); err != nil {
			log.Println("Fail to send message:", err)
		}
	}
}
