package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	apiURL  = "https://free.currconv.com/api/v7/convert"
	keysTTL = 5 * time.Minute
)

type Converter interface {
	ConvertRate(from, to string) (float64, error)
}

type CurrConvCom struct {
	httpClient *http.Client
	apiKey     string
	cache      map[string]cacheItem
}

type cacheItem struct {
	value   float64
	addTime time.Time
}

func New(apiKey string, timeout time.Duration) *CurrConvCom {
	conv := &CurrConvCom{
		httpClient: &http.Client{Timeout: timeout},
		apiKey:     apiKey,
		cache:      make(map[string]cacheItem),
	}
	return conv
}

func (converter CurrConvCom) ConvertRate(fromCurrency, toCurrency string) (float64, error) {
	var key = fromCurrency + "_" + toCurrency

	if res, ok := converter.cache[key]; ok {
		if time.Since(res.addTime) < keysTTL {
			return res.value, nil
		}
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 1.0, fmt.Errorf("fail to create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("q", key)
	q.Add("compact", "ultra")
	q.Add("apiKey", converter.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := converter.httpClient.Do(req)
	if err != nil {
		return 1.0, fmt.Errorf("fail to do request: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("fail to close request body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 1.0, fmt.Errorf("request failed: %v", err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 1.0, fmt.Errorf("fail to read request body: %v", err)
	}

	var parsed map[string]float64

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return 1.0, fmt.Errorf("fail to unmarshal request body: %v", err)
	}

	rate, ok := parsed[key]
	if !ok {
		return 1.0, fmt.Errorf("%s not found in %v", key, parsed)
	} else {
		converter.cache[key] = cacheItem{value: rate, addTime: time.Now()}
	}

	return rate, nil
}
