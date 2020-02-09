package investbot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"investbot/telegram"
	"investbot/tinkoff"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *telegram.Bot

func getEnvOrDie(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("ENV var %s is not set\n", key)
	}
	return val
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	telegramBotToken := getEnvOrDie("TELEGRAM_APITOKEN")
	tinkofApiToken := getEnvOrDie("TINKOFF_APITOKEN")
	webHookToken := getEnvOrDie("WEBHOOK_TOKEN")
	botOwnerId_ := getEnvOrDie("BOT_OWNER_ID")
	currencyConvertToken := getEnvOrDie("CURRENCY_API_TOKEN")

	botOwnerId, err := strconv.Atoi(botOwnerId_)
	if err != nil {
		log.Fatalln("Fail to parse bot owner ID:", botOwnerId_)
	}

	tinkoffApi := tinkoff.New(tinkofApiToken)

	if err := tinkoffApi.SetupAccounts(); err != nil {
		log.Fatalln("Fail to setup Tinkoff accounts:", err)
	}

	bot, err = telegram.New(telegramBotToken, webHookToken, currencyConvertToken, botOwnerId, tinkoffApi)
	if err != nil {
		log.Fatalln("Fail to create Telegram bot:", err)
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
		http.Error(w, "Fail to parse JSON body", http.StatusBadRequest)
		return
	}

	log.Printf("New update: %+v\n", update)

	if update.Message != nil {
		log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			bot.HandleCommandMessage(&update)
		}
	}

	w.WriteHeader(http.StatusOK)
}
