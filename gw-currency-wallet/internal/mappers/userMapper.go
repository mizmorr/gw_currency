package mappers

import (
	"errors"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/domain"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/store"
)

func ToStoreUserFromRegister(userReq *domain.RegisterRequest) (*store.User, error) {
	if userReq == nil {
		return nil, errors.New("user request is nil")
	}
	return &store.User{
		Username: userReq.Username,
		Email:    userReq.Email,
		Password: userReq.Password,
	}, nil
}

func ToStoreUserFromAuthorize(userReq *domain.AuthorizationRequst) (*store.User, error) {
	if userReq == nil {
		return nil, errors.New("user request is nil")
	}
	return &store.User{
		Username: userReq.Username,
		Password: userReq.Password,
	}, nil
}
