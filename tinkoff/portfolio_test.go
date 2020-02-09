package tinkoff

import (
	"encoding/json"
	"fmt"
	"testing"
)

const portfolioJson = `
{
  "trackingId": "c1672bac1e0c780d",
  "payload": {
    "positions": [
      {
        "figi": "BBG004S68829",
        "ticker": "TATNP",
        "isin": "RU0006944147",
        "instrumentType": "Stock",
        "balance": 72,
        "lots": 72,
        "expectedYield": {
          "currency": "RUB",
          "value": 295.2
        },
        "averagePositionPrice": {
          "currency": "RUB",
          "value": 718.5
        },
        "name": "Татнефть - привилегированные акции"
      },
      {
        "figi": "BBG00HYXLFQ7",
        "ticker": "RU000A0ZYEB1",
        "isin": "RU000A0ZYEB1",
        "instrumentType": "Bond",
        "balance": 52,
        "lots": 52,
        "expectedYield": {
          "currency": "RUB",
          "value": 405.3
        },
        "averagePositionPrice": {
          "currency": "RUB",
          "value": 1030.6
        },
        "averagePositionPriceNoNkd": {
          "currency": "RUB",
          "value": 1000.11
        },
        "name": "ТрансФин-М 001Р выпуск 4"
      },
      {
        "figi": "BBG005HLSZ23",
        "ticker": "FXUS",
        "isin": "IE00BD3QHZ91",
        "instrumentType": "Etf",
        "balance": 1,
        "lots": 1,
        "expectedYield": {
          "currency": "RUB",
          "value": 158
        },
        "averagePositionPrice": {
          "currency": "RUB",
          "value": 3660
        },
        "name": "Акции американских компаний"
      },
      {
        "figi": "BBG005HLTYH9",
        "ticker": "FXIT",
        "isin": "IE00BD3QJ757",
        "instrumentType": "Etf",
        "balance": 50,
        "lots": 50,
        "expectedYield": {
          "currency": "RUB",
          "value": 8350
        },
        "averagePositionPrice": {
          "currency": "RUB",
          "value": 6043
        },
        "name": "Акции компаний IT-сектора США"
      }
    ]
  },
  "status": "Ok"
}`

type mockedCurrencyConverter struct{}

func (c mockedCurrencyConverter) ConvertRate(from, to string) (float64, error) {
	return 1.5, nil
}

func TestPortfolio_Unmarshalling(t *testing.T) {
	var portfolio Portfolio
	err := json.Unmarshal([]byte(portfolioJson), &portfolio)
	if err != nil {
		t.Errorf("cannot unmarshal JSON to Portfolio struct: %v", err)
	}
	fmt.Printf("Portfolio:\n %+v", portfolio)
}

func TestPortfolio_Prettify(t *testing.T) {
	var portfolio Portfolio
	err := json.Unmarshal([]byte(portfolioJson), &portfolio)
	if err != nil {
		t.Errorf("fail to unmarshal portfolio: %v", err)
	}

	conv := mockedCurrencyConverter{}

	p, err := portfolio.Prettify(conv)
	if err != nil {
		t.Errorf("fail to prettify portfolio: %v\n", err)
	}
	fmt.Println(p)
}
