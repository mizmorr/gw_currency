package exchanger

type Rate struct {
	CurrencyCode string  `json:"currencyCode"`
	Value        float64 `json:"value"`
}
