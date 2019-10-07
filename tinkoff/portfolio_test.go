package tinkoff

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"text/template"
	"time"
	"tinkoff-invest-telegram-bot/currency"
)

const portfolioJson = `{
    "payload": {
        "positions": [
            {
                "balance": 300,
                "expectedYield": {
                    "currency": "RUB",
                    "value": 1314
                },
                "figi": "BBG0047315Y7",
                "instrumentType": "Stock",
                "isin": "RU0009029557",
                "lots": 30,
                "ticker": "SBERP"
            },
            {
                "balance": 4,
                "expectedYield": {
                    "currency": "USD",
                    "value": 14.87
                },
                "figi": "BBG000F1ZSQ2",
                "instrumentType": "Stock",
                "isin": "US57636Q1040",
                "lots": 4,
                "ticker": "MA"
            },
            {
                "balance": 3,
                "expectedYield": {
                    "currency": "USD",
                    "value": -10.8
                },
                "figi": "BBG000PSKYX7",
                "instrumentType": "Stock",
                "isin": "US92826C8394",
                "lots": 3,
                "ticker": "V"
            },
            {
                "balance": 200,
                "expectedYield": {
                    "currency": "RUB",
                    "value": -243.5
                },
                "figi": "BBG004S684M6",
                "instrumentType": "Stock",
                "isin": "RU0009062467",
                "lots": 20,
                "ticker": "SIBN"
            },
            {
                "balance": 3,
                "expectedYield": {
                    "currency": "RUB",
                    "value": 72
                },
                "figi": "BBG004731489",
                "instrumentType": "Stock",
                "isin": "RU0007288411",
                "lots": 3,
                "ticker": "GMKN"
            },
            {
                "balance": 55,
                "expectedYield": {
                    "currency": "RUB",
                    "value": -3534.4
                },
                "figi": "BBG006L8G4H1",
                "instrumentType": "Stock",
                "isin": "NL0009805522",
                "lots": 55,
                "ticker": "YNDX"
            },
            {
                "balance": 0.03,
                "expectedYield": {
                    "currency": "RUB",
                    "value": -0.02
                },
                "figi": "BBG0013HGFT4",
                "instrumentType": "Currency",
                "lots": 0,
                "ticker": "USD000UTSTOM"
            }
        ]
    },
    "status": "Ok",
    "trackingId": "8292e84e40c6f2bd"
}`

func TestPortfolio_Unmarshalling(t *testing.T) {
	var portfolio Portfolio
	err := json.Unmarshal([]byte(portfolioJson), &portfolio)
	if err != nil {
		t.Errorf("Cannot unmarshal JSON to Portfolio struct: %v", err)
	}
}

func TestPortfolio_Prettify(t *testing.T) {
	var portfolio Portfolio
	_ = json.Unmarshal([]byte(portfolioJson), &portfolio)

	tpl := template.Must(template.New("Portfolio").Funcs(PortfolioFuncMap).Parse(PortfolioTemplate))
	conv := currency.NewConverter(os.Getenv("CURRENCY_API_TOKEN"), 3*time.Second)

	p, err := portfolio.Prettify(tpl, conv)
	if err != nil {
		t.Errorf("Fail to prettify portfolio: %v\n", err)
	}
	fmt.Println(p)
}
