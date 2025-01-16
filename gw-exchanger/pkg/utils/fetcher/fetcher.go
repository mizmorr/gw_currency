package fetcher

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const apiURL = "https://api.exchangerate-api.com/v4/latest/RUB"

type apiResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

var currenciesName = []string{"RUB", "EUR", "USD"}

func FetchRates(ctx context.Context) (map[string]float64, error) {
	currencies := make(map[string]float64, 3)

	response, err := http.Get(apiURL)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка при выполнении запроса")
	}
	defer response.Body.Close()

	var data apiResponse
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "Ошибка при декодировании ответа")
	}
	for _, currency := range currenciesName {
		rate, ok := data.Rates[currency]
		if !ok {
			return nil, errors.Errorf("Курс для валюты %s не найден", currency)
		}
		currencies[currency] = rate
	}
	return currencies, nil
}
