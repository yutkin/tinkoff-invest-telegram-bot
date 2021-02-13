package tinkoff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Api struct {
	httpClient *http.Client
	token      string
	accounts   Accounts
}

func New(tinkoffAPIToken string) *Api {
	api := Api{
		token:      tinkoffAPIToken,
		httpClient: &http.Client{Timeout: 3 * time.Second},
	}

	return &api
}

func (api *Api) SetupAccounts() error {
	var err error
	api.accounts, err = api.GetAccounts()
	return err
}

func (api *Api) getPortfolio(accountId string) (Portfolio, error) {
	var portfolio Portfolio

	req, err := http.NewRequest("GET", PortfolioEndpoint, nil)
	if err != nil {
		return portfolio, fmt.Errorf("fail to create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("brokerAccountId", accountId)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", "Bearer "+api.token)
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return portfolio, fmt.Errorf("fail to do request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("fail to close request body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return portfolio, fmt.Errorf("fail to fetch portfolio. Status code: [%d]", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return portfolio, fmt.Errorf("fail to read response body: %v", err)
	}

	err = json.Unmarshal(respBody, &portfolio)

	if err != nil {
		return portfolio, fmt.Errorf("fail to unmarshal response body: %v", err)
	}

	return portfolio, nil
}

func (api *Api) GetIISPortfolio() (Portfolio, error) {
	if iisId, ok := api.accounts.IisAccountId(); ok {
		return api.getPortfolio(iisId)
	}
	return Portfolio{}, fmt.Errorf("iis account does not exist")
}

func (api *Api) GetBrokerPortfolio() (Portfolio, error) {
	if brokerId, ok := api.accounts.BrokerAccountId(); ok {
		return api.getPortfolio(brokerId)
	}
	return Portfolio{}, fmt.Errorf("broker account does not exist")
}
