package tinkoff

import (
	"bytes"
	"fmt"
	"text/template"
)

const PortfolioTemplate = `
{{- range $i, $v := .Payload.Positions}}
	{{- inc $i}}. <b>{{.Ticker}}</b> {{.Balance}} ({{sign .ExpectedYield.Value}} {{.ExpectedYield.Currency}})
{{end}}`

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