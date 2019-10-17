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
                "averagePositionPrice": {
                    "currency": "RUB",
                    "value": 194.24
                },
                "balance": 250,
                "expectedYield": {
                    "currency": "RUB",
                    "value": 2263
                },
                "figi": "BBG0047315Y7",
                "instrumentType": "Stock",
                "isin": "RU0009029557",
                "lots": 25,
                "ticker": "SBERP"
            },
            {
                "averagePositionPrice": {
                    "currency": "USD",
                    "value": 270.08
                },
                "balance": 4,
                "expectedYield": {
                    "currency": "USD",
                    "value": 32.23
                },
                "figi": "BBG000F1ZSQ2",
                "instrumentType": "Stock",
                "isin": "US57636Q1040",
                "lots": 4,
                "ticker": "MA"
            },
            {
                "averagePositionPrice": {
                    "currency": "USD",
                    "value": 179.5
                },
                "balance": 3,
                "expectedYield": {
                    "currency": "USD",
                    "value": -5.43
                },
                "figi": "BBG000PSKYX7",
                "instrumentType": "Stock",
                "isin": "US92826C8394",
                "lots": 3,
                "ticker": "V"
            },
            {
                "averagePositionPrice": {
                    "currency": "RUB",
                    "value": 412.55
                },
                "balance": 130,
                "expectedYield": {
                    "currency": "RUB",
                    "value": 2036.5
                },
                "figi": "BBG004S684M6",
                "instrumentType": "Stock",
                "isin": "RU0009062467",
                "lots": 13,
                "ticker": "SIBN"
            },
            {
                "averagePositionPrice": {
                    "currency": "RUB",
                    "value": 15870
                },
                "balance": 3,
                "expectedYield": {
                    "currency": "RUB",
                    "value": 858
                },
                "figi": "BBG004731489",
                "instrumentType": "Stock",
                "isin": "RU0007288411",
                "lots": 3,
                "ticker": "GMKN"
            },
            {
                "averagePositionPrice": {
                    "currency": "RUB",
                    "value": 5717.5
                },
                "balance": 5,
                "expectedYield": {
                    "currency": "RUB",
                    "value": -37
                },
                "figi": "BBG004731032",
                "instrumentType": "Stock",
                "isin": "RU0009024277",
                "lots": 5,
                "ticker": "LKOH"
            },
            {
                "averagePositionPrice": {
                    "currency": "RUB",
                    "value": 2160.4
                },
                "balance": 37,
                "expectedYield": {
                    "currency": "RUB",
                    "value": -7521.2
                },
                "figi": "BBG006L8G4H1",
                "instrumentType": "Stock",
                "isin": "NL0009805522",
                "lots": 37,
                "ticker": "YNDX"
            },
            {
                "averagePositionPrice": {
                    "currency": "RUB",
                    "value": 65.44
                },
                "balance": 0.03,
                "expectedYield": {
                    "currency": "RUB",
                    "value": -0.04
                },
                "figi": "BBG0013HGFT4",
                "instrumentType": "Currency",
                "lots": 0,
                "ticker": "USD000UTSTOM"
            }
        ]
    },
    "status": "Ok",
    "trackingId": "b501c8e61e7e8f7d"
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
