package storage

import (
	"context"

	"github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/model"
)

type Repository interface {
	GetAllRates(ctx context.Context) ([]*model.Rate, error)
	GetRate(ctx context.Context, code string) (*model.Rate, error)

	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
