package tinkoff

import (
	"bytes"
	"html/template"
)

const PortfolioTemplate = "{{range $i, $v := .Payload.Positions}}" +
	"{{inc $i}}. <b>{{.Ticker}}</b> {{.Balance}} ({{.ExpectedYield.Value}} {{.ExpectedYield.Currency}})\n" +
	"{{end}}"

var PortfolioFuncMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
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
			ExpectedYield  struct {
				Currency string  `json:"currency"`
				Value    float64 `json:"value"`
			}
			Lots int32 `json:"lots"`
		}
	} `json:"payload"`
}

func (portfolio *Portfolio) Prettify(t *template.Template) (string, error) {
	buff := bytes.Buffer{}
	err := t.Execute(&buff, portfolio)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}