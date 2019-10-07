package tinkoff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"tinkoff-invest-telegram-bot/currency"
)

const (
	ApiUrl = "https://api-invest.tinkoff.ru/openapi/portfolio"
)

type Api struct {
	Token             string
	Client            *http.Client
	PortfolioTemplate *template.Template
	CurrencyConverter *currency.Converter
}

func (api *Api) GetPortfolio() (Portfolio, error) {
	req, err := http.NewRequest("GET", ApiUrl, nil)
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
		return Portfolio{}, fmt.Errorf("Fail to fetch portfolio. Status code: [%d]", resp.StatusCode)
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
