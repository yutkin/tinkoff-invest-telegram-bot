package tinkoff_investments_telegram_bot

import (
	"encoding/json"
	"text/template"
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
	webHookToken := os.Getenv("WEBHOOK_TOKEN")
	botOwnerId_ := os.Getenv("BOT_OWNER_ID")

	botOwnerId, err := strconv.Atoi(botOwnerId_)

	if err != nil {
		log.Fatalf("Cannot parse bot owner ID: %s", botOwnerId_)
	}

	api, err := tgbotapi.NewBotAPI(telegramBotToken)

	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", api.Self.UserName)

	bot.TelegramgApi = api
	bot.OwnerId = botOwnerId
	bot.WebHookToken = webHookToken

	bot.TinkoffApi = &tinkoff.Api{
		Url:    tinkoff.URL,
		Token:  tinkofApiToken,
		Client: &http.Client{Timeout: tinkoff.TIMEOUT},
		PortfolioTemplate: template.Must(
			template.New("Portfolio").Funcs(tinkoff.PortfolioFuncMap).Parse(tinkoff.PortfolioTemplate),
		),
	}
}

func HandleTelegramUpdate(w http.ResponseWriter, r *http.Request) {

	webHookToken := r.URL.Query().Get("token")

	if webHookToken != bot.WebHookToken {
		http.Error(w, "Bad token", http.StatusUnauthorized)
		return
	}

	update := tgbotapi.Update{}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Can't parse JSON body", http.StatusBadRequest)
		return
	}

	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			bot.HandleCommandMessage(&update)
		}
	}

	w.WriteHeader(http.StatusOK)
}
