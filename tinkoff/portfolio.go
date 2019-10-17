package tinkoff

import (
	"bytes"
	"fmt"
	"text/template"
	"tinkoff-invest-telegram-bot/currency"
)

const (
	PortfolioTemplate = `
{{- range $i, $v := .Payload.Positions}}
	{{- inc $i}}. <b>{{.Name}}</b> {{.Balance}} {{if ne .InstrumentType "Currency"}}шт. {{end}} {{formatFloat .TotalPositionPrice}} {{.AveragePositionPrice.Currency}} ({{sign .ExpectedYield.Value}} {{.ExpectedYield.Currency}})
{{end}}
Итог: <b>{{sign .TotalYieldRUB}}</b> RUB
`
)

var PortfolioFuncMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"sign": func(i interface{}) string {
		switch v := i.(type) {
		case float32, float64:
			return fmt.Sprintf("%+.2f", v)
		case int, int8, int32, int64:
			return fmt.Sprintf("%+d", v)
		default:
			return ""
		}
	},
	"formatFloat": func(i interface{}) string {
		switch v := i.(type) {
		case float32, float64:
			return fmt.Sprintf("%.2f", v)
		default:
			return ""
		}
	},
}

type Portfolio struct {
	TrackingID string `json:"trackingId"`
	Status     string `json:"status"`
	Payload    struct {
		Positions []struct {
			Figi           string  `json:"figi"`
			Ticker         string  `json:"ticker"`
			Balance        float64 `json:"balance"`
			InstrumentType string  `json:"instrumentType"`
			Lots           int32   `json:"lots"`

			ExpectedYield struct {
				Currency string  `json:"currency"`
				Value    float64 `json:"value"`
			}

			AveragePositionPrice struct {
				Currency string  `json:"currency"`
				Value    float64 `json:"value"`
			}
			Name               string
			TotalPositionPrice float64
		}
	} `json:"payload"`

	TotalYieldRUB float64
}

func (portfolio *Portfolio) Prettify(t *template.Template, converter *currency.Converter) (string, error) {

	for i, v := range portfolio.Payload.Positions {
		var rate float64 = 1.0

		if v.ExpectedYield.Currency != "RUB" {
			newRate, err := converter.GetCurrencyConvertRate(v.ExpectedYield.Currency, "RUB")
			if err != nil {
				return "", fmt.Errorf("Fail to fetch %s to RUB exchange rate: %v", v.ExpectedYield.Currency, err)
			}
			rate = newRate
		}

		portfolio.TotalYieldRUB += v.ExpectedYield.Value * rate
		portfolio.Payload.Positions[i].TotalPositionPrice = v.Balance*v.AveragePositionPrice.Value + v.ExpectedYield.Value
	}

	buff := bytes.Buffer{}
	err := t.Execute(&buff, portfolio)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}
