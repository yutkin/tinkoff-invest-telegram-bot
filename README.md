# Деплой

```
gcloud functions deploy HandleTelegramUpdate \
    --runtime go111 \
    --trigger-http \
    --memory=128 \
    --timeout=10 \
    --region=europe-west2 \
    --set-env-vars=TELEGRAM_APITOKEN=$TELEGRAM_APITOKEN,TINKOFF_APITOKEN=$TINKOFF_APITOKEN,BOT_OWNER_ID=$BOT_OWNER_ID
```
    
# Установка Web-Хука    
`http -v https://api.telegram.org/bot$TELEGRAM_APITOKEN/setWebhook url="https://europe-west2-caramel-sum-255011.cloudfunctions.net/HandleTelegramUpdate "`