package store

import (
	"time"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Wallet struct {
	ID        int64
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateBalance struct {
	UserID    int64
	Currency  string
	Amount    float64
	Operation string
}

type WalletCurrency struct {
	ID        int64
	WalletID  int64
	Currency  string
	Balance   float64
	UpdatedAt time.Time
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	Hash      string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

type ExchangeBalance struct {
	UserID       int64
	FromCurrency string
	FromAmount   float64
	ToCurrency   string
	ToAmount     float64
}

type CurrencyRequest struct {
	UserID       int64
	CurrencyCode string
}
