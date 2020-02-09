package tinkoff

import (
	"bytes"
	"fmt"
	"text/template"
	"investbot/currency"
)

const PortfolioEndpoint = "https://api-invest.tinkoff.ru/openapi/portfolio"

const (
	portfolioTemplateStr = `
{{- range $i, $v := .Payload.Positions}}
<u>{{.Name}}</u> <b>({{sign .ExpectedYield.Value}}</b> {{.ExpectedYield.Currency}})
Баланс: {{.Balance}} {{if ne .InstrumentType "Currency"}}шт. на{{end}} {{formatFloat .TotalPositionPrice}} {{.AveragePositionPrice.Currency}}
{{end}}
<b>Итог: {{sign .TotalYieldRUB}} RUB</b>
`
)

var portfolioTemplate = template.Must(
	template.New("Portfolio").Funcs(portfolioFuncMap).Parse(portfolioTemplateStr),
)

var portfolioFuncMap = template.FuncMap{
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
	TotalYieldRUB float64

	TrackingID string `json:"trackingId"`
	Status     string `json:"status"`

	Payload struct {
		Positions []struct {
			Figi           string  `json:"figi"`
			Ticker         string  `json:"ticker"`
			Balance        float64 `json:"balance"`
			InstrumentType string  `json:"instrumentType"`
			Lots           int32   `json:"lots"`
			Blocked        float64 `json:"blocked"`

			ExpectedYield struct {
				Currency string  `json:"currency"`
				Value    float64 `json:"value"`
			}

			AveragePositionPrice struct {
				Currency string  `json:"currency"`
				Value    float64 `json:"value"`
			}

			AveragePositionPriceNoNkd struct {
				Currency string  `json:"currency"`
				Value    float64 `json:"value"`
			}

			Name               string `json:"name"`
			TotalPositionPrice float64
		}
	} `json:"payload"`
}

func (portfolio *Portfolio) Prettify(converter currency.Converter) (string, error) {

	for i, v := range portfolio.Payload.Positions {
		rate := 1.0

		if v.ExpectedYield.Currency != "RUB" {
			newRate, err := converter.ConvertRate(v.ExpectedYield.Currency, "RUB")
			if err != nil {
				newErr := fmt.Errorf("fail to fetch %s to RUB exchange rate: %v", v.ExpectedYield.Currency, err)
				return "", newErr
			}
			rate = newRate
		}

		portfolio.TotalYieldRUB += v.ExpectedYield.Value * rate
		portfolio.Payload.Positions[i].TotalPositionPrice = v.Balance*v.AveragePositionPrice.Value + v.ExpectedYield.Value
	}

	buff := bytes.Buffer{}

	err := portfolioTemplate.Execute(&buff, portfolio)
	if err != nil {
		return "", fmt.Errorf("fail to execute tempalte: %v", err)
	}

	return buff.String(), nil
}
