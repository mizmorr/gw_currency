package domain

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthorizationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type BalanceResponse struct {
	Currency string  `json:"currency" `
	Value    float64 `json:"value" `
}

type DepositRequest struct {
	Currency string  `json:"currency" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
}

type WithdrawRequest struct {
	Currency string  `json:"currency" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
}

type RateResponse struct {
	CurrencyCode string  `json:"currency_code" binding:"required"`
	Value        float64 `json:"value" binding:"required"`
}

type ExchangeRequest struct {
	BaseCurrency   string  `json:"base_currency" binding:"required"`
	TargetCurrency string  `json:"target_currency" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
}

type ExchangeResponse struct {
	ExchangeAmount float64            `json:"exchange_amount"`
	Message        string             `json:"message"`
	NewBalance     []*BalanceResponse `json:"new_balance"`
}

type TokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type RefreshRequest struct {
	TokenHash string `json:"tokenhash" binding:"required"`
}
