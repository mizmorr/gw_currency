package pg

import (
	"context"
	"fmt"
	"testing"

	"github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/model"
	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	ctx := context.Background()

	repo, err := NewPostgresRepo(ctx)
	assert.Nil(t, err)
	err = repo.Start(ctx)

	assert.Nil(t, err)

	newrate := &model.Rate{
		CurrencyCode: "RUB",
		Value:        1,
	}

	rateFromDB, err := repo.GetRate(ctx, "RUB")
	fmt.Println(rateFromDB)
	assert.Nil(t, err)

	assert.Equal(t, newrate, rateFromDB)
}

func TestGetAllRates(t *testing.T) {
	ctx := context.Background()

	repo, err := NewPostgresRepo(ctx)
	assert.Nil(t, err)
	err = repo.Start(ctx)

	assert.Nil(t, err)

	rates := []*model.Rate{
		{
			CurrencyCode: "USD",
		},
		{
			CurrencyCode: "RUB",
		},

		{
			CurrencyCode: "EUR",
		},
	}

	allRates, err := repo.GetAllRates(ctx)
	assert.Nil(t, err)
	assert.Len(t, allRates, 3)
	for i, r := range allRates {
		assert.Equal(t, rates[i].CurrencyCode, r.CurrencyCode)
		assert.Positive(t, r.Value)
	}
}
