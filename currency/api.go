package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	CurrencyApiUrl = "https://free.currconv.com/api/v7/convert"
	fiveMinutes    = 5 * time.Minute
)

type cacheItem struct {
	value   float64
	addTime time.Time
}

type Converter struct {
	Api    *http.Client
	ApiKey string
	cache  map[string]cacheItem
}

func NewConverter(apiKey string, timeout time.Duration) *Converter {
	conv := &Converter{
		Api:    &http.Client{Timeout: timeout},
		ApiKey: apiKey,
		cache:  make(map[string]cacheItem),
	}
	return conv
}

func (converter *Converter) GetCurrencyConvertRate(fromCurrency, toCurrency string) (float64, error) {
	var key = fromCurrency + "_" + toCurrency

	if res, ok := converter.cache[key]; ok {
		if time.Now().Sub(res.addTime) < fiveMinutes {
			return res.value, nil
		}
	}

	req, err := http.NewRequest("GET", CurrencyApiUrl, nil)
	if err != nil {
		return 1.0, err
	}

	q := req.URL.Query()
	q.Add("q", key)
	q.Add("compact", "ultra")
	q.Add("apiKey", converter.ApiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := converter.Api.Do(req)
	if err != nil {
		return 1.0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return 1.0, fmt.Errorf("Fail to fetch portfolio: [%d] %s", resp.StatusCode, string(b))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 1.0, err
	}

	var parsed map[string]float64

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return 1.0, err
	}

	rate, ok := parsed[key]
	if !ok {
		return 1.0, fmt.Errorf("%s not found in %v", key, parsed)
	} else {
		converter.cache[key] = cacheItem{value: rate, addTime: time.Now()}
	}

	return rate, nil
}
