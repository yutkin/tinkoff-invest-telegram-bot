package tinkoff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"tinkoff-invest-telegram-bot/currency"
)

const (
	ApiUrl = "https://api-invest.tinkoff.ru/openapi"
)

type Api struct {
	Token             string
	Client            *http.Client
	PortfolioTemplate *template.Template
	CurrencyConverter *currency.Converter
	FigiToName        map[string]string
}

func (api *Api) FillPortfolioPositionsNames(portfolio *Portfolio) {

	api.FigiToName = make(map[string]string)

	type FigiName struct {
		Figi string
		Name string
	}

	type resultChanel chan FigiName

	fetchPositionName := func(figi string, resultChan resultChanel) {
		log.Println("Trying to GET Name for", figi)
		req, err := http.NewRequest("GET", ApiUrl+"/market/search/by-figi", nil)
		if err != nil {
			resultChan <- FigiName{Figi: figi}
			return
		}

		q := req.URL.Query()
		q.Add("figi", figi)

		req.URL.RawQuery = q.Encode()
		req.Header.Add("Authorization", "Bearer "+api.Token)

		resp, err := api.Client.Do(req)
		if err != nil {
			resultChan <- FigiName{Figi: figi}
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			resultChan <- FigiName{Figi: figi}
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChan <- FigiName{Figi: figi}
			return
		}

		res := struct {
			Payload struct {
				Name string `json:"name"`
			}
		}{}

		err = json.Unmarshal(respBody, &res)

		if err != nil {
			resultChan <- FigiName{Figi: figi}
			return
		}

		resultChan <- FigiName{Figi: figi, Name: res.Payload.Name}
	}

	ch := make(resultChanel)

	var n int = 0
	for _, v := range portfolio.Payload.Positions {
		if _, ok := api.FigiToName[v.Figi]; !ok {
			go fetchPositionName(v.Figi, ch)
			n++
		}
	}

	for i := 0; i < n; i++ {
		v := <- ch
		api.FigiToName[v.Figi] = v.Name
		log.Println("Result:", v)
	}

	for i, _ := range portfolio.Payload.Positions {
		figi := portfolio.Payload.Positions[i].Figi
		portfolio.Payload.Positions[i].Name = api.FigiToName[figi]
	}
	log.Println("FillPortfolioPositionsNames done.")
}

func (api *Api) GetPortfolio() (Portfolio, error) {
	req, err := http.NewRequest("GET", ApiUrl+"/portfolio", nil)
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
