package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JacekKaczmarczyk/Rates-API/utils"
)

type NbpResponse struct {
	Table         string `json:"table"`
	No            string `json:"no"`
	EffectiveDate string `json:"effectiveDate"`
	Rates         []Rate `json:"rates"`
}

type Rate struct {
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Mid      float64 `json:"mid"`
}

type NbpProvider struct {
	Name       string
	APIURL     string
	DateFormat string
}

func NewNbpProvider() *NbpProvider {
	return &NbpProvider{
		Name:       "NBP",
		APIURL:     "https://api.nbp.pl/api/exchangerates/tables/a",
		DateFormat: "2006-01-02",
	}
}

func (p *NbpProvider) GetCurrencies(codes []string, date string) (Response, error, int) {
	// Validate currency code formats before making request
	for _, code := range codes {
		if !utils.ValidateCurrencyCodeFormat(code) {
			return Response{}, fmt.Errorf("invalid currency code format: %s", code), http.StatusBadRequest
		}
	}

	req, err := p.createGetRequest(date)
	if err != nil {
		return Response{}, err, http.StatusBadRequest
	}

	response, err, statusCode := p.fetchNbpData(req)
	if err != nil {
		return Response{}, err, statusCode
	}

	filteredRates := p.filterRates(response[0].Rates, codes)
	if len(filteredRates) == 0 {
		return Response{}, fmt.Errorf("no rates found for the specified codes: %v", codes), http.StatusNotFound
	}

	return Response{
		AsOF:     response[0].EffectiveDate,
		Provider: p.Name,
		Rates:    filteredRates,
	}, nil, http.StatusOK
}

func (p *NbpProvider) createGetRequest(date string) (*http.Request, error) {
	uri := p.APIURL
	if date != "" {
		if !utils.ValidateDate(date, p.DateFormat) {
			return nil, fmt.Errorf("invalid date format: %s, expected format: %s", date, p.DateFormat)
		}
		uri = fmt.Sprintf("%s/%s", p.APIURL, date)
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Go-NBP-Client/1.0")
	return req, nil
}

func (p *NbpProvider) fetchNbpData(req *http.Request) ([]NbpResponse, error, int) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode >= 500 {
			return nil, fmt.Errorf("NBP API unavailable (status %d): %s", res.StatusCode, string(body)), http.StatusBadGateway
		}
		return nil, fmt.Errorf("NBP API returned status %d: %s", res.StatusCode, string(body)), http.StatusBadRequest
	}

	var response []NbpResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return response, nil, http.StatusOK
}

func (p *NbpProvider) filterRates(rates []Rate, codes []string) []RateValue {
	codeMap := make(map[string]bool)
	for _, code := range codes {
		codeMap[code] = true
	}

	result := make([]RateValue, 0)
	for _, rate := range rates {
		if codeMap[rate.Code] {
			result = append(result, RateValue{
				Code:  rate.Code,
				Value: rate.Mid,
			})
		}
	}

	return result
}
