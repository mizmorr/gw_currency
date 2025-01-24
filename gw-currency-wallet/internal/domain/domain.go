package domain

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRespone struct {
	Message string `json:"message"`
}

type AuthorizationRequst struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthorizationResponse struct {
	Token string `json:"token"`
}

type BalanceResponse struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type DepositRequest struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type WithdrawRequest struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type RateResponse struct {
	CurrencyCode string  `json:"currency_code"`
	Value        float64 `json:"value"`
}

type ExchangeRequest struct {
	BaseCurrency   string  `json:"base_currency"`
	TargetCurrency string  `json:"target_currency"`
	Amount         float64 `json:"amount"`
}

type ExchangeResponse struct {
	ExchangeAmount float64            `json:"exchange_amount"`
	Message        string             `json:"message"`
	NewBalance     []*BalanceResponse `json:"new_balance"`
}

type TokenRepsonse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type RefreshRequest struct {
	TokenHash string `json:"tokenhash"`
}
