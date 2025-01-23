package mappers

import (
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
)

func ToDomainBalance(storeBalance []*store.WalletCurrency) []*domain.BalanceResponse {
	balances := make([]*domain.BalanceResponse, 0, len(storeBalance))
	for _, b := range storeBalance {
		balance := &domain.BalanceResponse{
			Currency: b.Currency,
			Value:    b.Balance,
		}
		balances = append(balances, balance)
	}
	return balances
}

func ToStoreDepositBalance(userid int64, updateRequest *domain.DepositRequest) *store.UpdateBalance {
	return &store.UpdateBalance{
		UserID:    userid,
		Operation: "deposit",
		Currency:  updateRequest.Currency,
		Amount:    updateRequest.Amount,
	}
}

func ToStoreWithdrawBalance(userid int64, updateRequest *domain.WithdrawRequest) *store.UpdateBalance {
	return &store.UpdateBalance{
		Operation: "withdraw",
		UserID:    userid,
		Currency:  updateRequest.Currency,
		Amount:    updateRequest.Amount,
	}
}
