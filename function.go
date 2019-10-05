package tinkoff_investments_telegram_bot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"tinkoff-investments-telegram-bot/tgbot"
	"tinkoff-investments-telegram-bot/tinkoff"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot = tgbot.TinkoffInvestmentsBot{}

func init() {
	telegramBotToken := os.Getenv("TELEGRAM_APITOKEN")
	tinkofApiToken := os.Getenv("TINKOFF_APITOKEN")

	botOwnerId, err := strconv.Atoi(os.Getenv("BOT_OWNER_ID"))

	if err != nil {
		log.Fatalf("Cannot parse owner ID: %s", botOwnerId)
	}

	api, err := tgbotapi.NewBotAPI(telegramBotToken)

	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", api.Self.UserName)

	bot.TelegramgApi = api
	bot.OwnerId = botOwnerId

	bot.TinkoffApi = &tinkoff.Api{tinkoff.URL, tinkofApiToken, &http.Client{Timeout: tinkoff.TIMEOUT}}
}

func HandleTelegramUpdate(w http.ResponseWriter, r *http.Request) {

	update := tgbotapi.Update{}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Can't parse JSON body", http.StatusBadRequest)
		return
	}

	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() && update.Message.From.ID == bot.OwnerId {
			bot.HandleCommandMessage(&update)
		}
	}

	w.WriteHeader(http.StatusOK)
}
