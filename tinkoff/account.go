package tinkoff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const AccountsEndpoint = "https://api-invest.tinkoff.ru/openapi/user/accounts"

type Accounts struct {
	Payload struct {
		Accounts []struct {
			BrokerAccountType string `json:"brokerAccountType"`
			BrokerAccountId   string `json:"brokerAccountId"`
		} `json:"accounts"`
	} `json:"payload"`
}

func (accounts *Accounts) IisAccountId() (string, bool) {
	for _, acc := range accounts.Payload.Accounts {
		if acc.BrokerAccountType == "TinkoffIis" && acc.BrokerAccountId != "" {
			return acc.BrokerAccountId, true
		}
	}
	return "", false
}

func (accounts *Accounts) BrokerAccountId() (string, bool) {
	for _, acc := range accounts.Payload.Accounts {
		if acc.BrokerAccountType == "Tinkoff" && acc.BrokerAccountId != "" {
			return acc.BrokerAccountId, true
		}
	}
	return "", false
}

func (api *Api) GetAccounts() (Accounts, error) {
	var accounts Accounts

	req, err := http.NewRequest("GET", AccountsEndpoint, nil)

	if err != nil {
		return accounts, fmt.Errorf("fail to create request object: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+api.token)

	resp, err := api.httpClient.Do(req)

	if err != nil {
		return accounts, fmt.Errorf("fail to do request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("fail to close request body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return accounts, fmt.Errorf("fail to fetch accounts. Status code: [%d]", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return accounts, fmt.Errorf("fail to read body data: %v", err)
	}

	err = json.Unmarshal(respBody, &accounts)

	if err != nil {
		return accounts, fmt.Errorf("fail to unmarshal request body: %v", err)
	}

	return accounts, nil
}
