package tinkoff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	TIMEOUT = time.Second * 3
	URL     = "https://api-invest.tinkoff.ru/openapi/portfolio"
)

type Api struct {
	Url    string
	Token  string
	Client *http.Client
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

func (portfolio *Portfolio) Prettify() string {
	positions := make([]string, len(portfolio.Payload.Positions))

	for i, position := range portfolio.Payload.Positions {
		expectedYield := position.ExpectedYield.Value
		currency := position.ExpectedYield.Currency

		switch position.InstrumentType {
		case "Stock":
			positions = append(
				positions,
				fmt.Sprintf(
					"<b>%d.</b> %s %d (%+.2f %s)",
					i+1,
					position.Ticker,
					int(position.Balance),
					expectedYield,
					currency,
				),
			)
			//case "Currency":
			//	positions = append(
			//		positions,
			//		fmt.Sprintf("%15s %9.2f (%+.2f %s)", position.Ticker, position.Balance, expectedYield, currency),
			//	)
		}
	}

	return strings.Join(positions, "\n")
}

func (api *Api) GetPortfolio() (Portfolio, error) {
	req, err := http.NewRequest("GET", api.Url, nil)
	if err != nil {
		return Portfolio{}, err
	}

	req.Header.Add("Authorization", "Bearer "+api.Token)
	resp, err := api.Client.Do(req)
	if err != nil {
		return Portfolio{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Portfolio{}, fmt.Errorf("Fail to make request:", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Portfolio{}, err
	}

	var portfolio Portfolio
	err = json.Unmarshal(respBody, &portfolio)

	if err != nil {
		return Portfolio{}, err
	}

	return portfolio, nil
}
