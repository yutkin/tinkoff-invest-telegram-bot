# Описание
Telegram бот для получения портфеля в "Тинькофф Инвестиции".

Бот разработан под [serverless](https://ru.wikipedia.org/wiki/%D0%91%D0%B5%D1%81%D1%81%D0%B5%D1%80%D0%B2%D0%B5%D1%80%D0%BD%D1%8B%D0%B5_%D0%B2%D1%8B%D1%87%D0%B8%D1%81%D0%BB%D0%B5%D0%BD%D0%B8%D1%8F) деплой на [Google Cloud Functions](https://cloud.google.com/functions/). Serverless позволяет хостить бота бесплатно (или почти бесплатно) при небольших нагрузках на backend.

# Структура
* `tgbot` - модуль с обработчиком обновлений от Telegram бота
* `tinkoff` - модуль для работы с Tinkoff OpenAPI
* `function.go` - точка входа для Google Functions. Инициализация глобальных переменных и обработка запросов на WebHook.  

# Подготовка к деплою
1. Создаём Telegram бота и получаем для него токен. [Инструкция](https://core.telegram.org/bots#6-botfather).
2. Получаем токен в [Tinkoff Open API](https://tinkoffcreditsystems.github.io/invest-openapi/auth/).
3. Узнаём свой Telegram ID, например, через [@userinfobot](https://t.me/userinfobot).

# Деплой на Google Functions
1. Клонируем репозиторий `git clone https://github.com/yutkin/tinkoff-invest-telegram-bot.git && cd tinkoff-invest-telegram-bot`
2. Собираем зависимости: `go mod vendor`
3. Выставляем переменные среды окружения:
```
export TELEGRAM_APITOKEN=<telegram_bot_token>
export TINKOFF_APITOKEN=<tinkoff_api_token>
export BOT_OWNER_ID=<your_telegram_id>
export WEBHOOK_TOKEN=<any_random_string>
```
5. С помощью утилиты [gcloud](https://cloud.google.com/sdk/gcloud/) деплоим бота:
```
gcloud functions deploy HandleTelegramUpdate \
    --runtime go111 \
    --trigger-http \
    --memory=128 \
    --timeout=10 \
    --region=europe-west2 \
    --set-env-vars=TELEGRAM_APITOKEN=$TELEGRAM_APITOKEN,TINKOFF_APITOKEN=$TINKOFF_APITOKEN,BOT_OWNER_ID=$BOT_OWNER_ID,WEBHOOK_TOKEN=$WEBHOOK_TOKEN
```
    
[Описание флагов](https://cloud.google.com/sdk/gcloud/reference/functions/deploy) для команды `gcloud functions deploy`.

# Настройка Telegram Webhook
При успешном завершении, `gcloud functions deploy` печатает результат в `YAML` формате. Нас интересует поле `httpsTrigger.url`.
Значение `httpsTrigger.url` – endpoint нашего обработчика. Его и нужно использовать в качестве адреса WebHook. 

Установка WebHook происходит через метод [setWebhook](https://core.telegram.org/bots/api#setwebhook). Пример:
```
http -v https://api.telegram.org/bot$TELEGRAM_APITOKEN/setWebhook \
    url="https://europe-west2-<your_project>.cloudfunctions.net/HandleTelegramUpdate?token=$WEBHOOK_TOKEN"
```

После установки WebHook, Telegram будет отправлять все обновления от нашего бота в нашу функцию на Google Functions.

# Дополнительно
1. [Документация к Google Functions](https://cloud.google.com/functions/docs/)
2. [Документация к Telegram Bot API](https://core.telegram.org/bots/api)
3. [Репозиторий с документацией к Tinkoff OpenAPI](https://github.com/TinkoffCreditSystems/invest-openapi/)
4. [Golang telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)