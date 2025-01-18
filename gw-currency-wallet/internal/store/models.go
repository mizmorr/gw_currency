package store

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Wallet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WalletBalance struct {
	ID        uuid.UUID
	WalletID  uuid.UUID
	Currency  string
	Balance   float64
	UpdatedAt time.Time
}

type Transaction struct {
	ID          uuid.UUID
	WalletID    uuid.UUID
	Currency    string
	Amount      float64
	Type        string
	Description string
	CreatedAt   time.Time
}
