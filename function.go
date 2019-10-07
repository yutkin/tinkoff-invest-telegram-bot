package tinkoff_invest_telegram_bot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
	"tinkoff-invest-telegram-bot/currency"
	"tinkoff-invest-telegram-bot/tgbot"
	"tinkoff-invest-telegram-bot/tinkoff"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot = tgbot.TinkoffInvestmentsBot{}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
		Token:  tinkofApiToken,
		Client: &http.Client{Timeout: 5 * time.Second},
		PortfolioTemplate: template.Must(
			template.New("Portfolio").Funcs(tinkoff.PortfolioFuncMap).Parse(tinkoff.PortfolioTemplate),
		),
		CurrencyConverter: currency.NewConverter(os.Getenv("CURRENCY_API_TOKEN"), 5*time.Second),
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
