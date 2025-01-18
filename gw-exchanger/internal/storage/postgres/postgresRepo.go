package pg

import (
	"context"
	"time"

	"github.com/mizmorr/gw_currency/gw-exchanger/internal/config"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/storage"
	"github.com/mizmorr/gw_currency/gw-exchanger/internal/storage/model"
	"github.com/mizmorr/gw_currency/gw-exchanger/pkg/utils/fetcher"
	logger "github.com/mizmorr/loggerm"
	"github.com/pkg/errors"
)

type PostgresRepo struct {
	db     *pg
	stop   chan interface{}
	config *config.Config
	log    *logger.Logger
}

func NewPostgresRepo(ctx context.Context) (storage.Repository, error) {
	db, err := newPg(ctx)
	if err != nil {
		return nil, err
	}

	log := logger.GetLoggerFromContext(ctx)

	ch := make(chan interface{})

	return &PostgresRepo{
		db:     db,
		stop:   ch,
		config: config.Get(),
		log:    log,
	}, nil
}

func (repo *PostgresRepo) Start(ctx context.Context) error {
	err := dial(ctx)
	if err != nil {
		return err
	}
	go repo.keepAlive(ctx)

	go repo.updater(ctx)

	return nil
}

func (repo *PostgresRepo) keepAlive(ctx context.Context) {
	repo.log.Debug().Msg("Keeping database connection alive...")

	for {
		select {
		case <-repo.stop:
			repo.log.Info().Msg("Keep alive worker is stopped..")
			return
		default:
			repo.maintainConnection(ctx)
		}
	}
}

func (repo *PostgresRepo) maintainConnection(ctx context.Context) {
	time.Sleep(repo.config.KeepAliveTimeout)

	connectionLost := false

	conn, err := repo.db.Acquire(ctx)
	if err != nil {
		connectionLost = true
		repo.log.Debug().Msg("[keepAlive] Lost connection, is trying to reconnect...")
	} else {
		conn.Release()
	}

	if connectionLost {
		err = dial(ctx)
		if err != nil {
			repo.log.Err(err).Msg("Failed to reconnect to PostgreSQL database")
		}
	}
}

const workersCount = 2

func (repo *PostgresRepo) Stop(_ context.Context) error {
	repo.log.Info().Msg("Stopping PostgreSQL repository..")
	for range workersCount {
		repo.stop <- struct{}{}
	}
	repo.db.Close()
	return nil
}

func (repo *PostgresRepo) updater(ctx context.Context) {
	repo.update(ctx)
	updateStamp := time.Now()

	repo.log.Info().Msg("Updater worker is starting..")

	for {
		select {
		case <-repo.stop:
			repo.log.Info().Msg("Updater worker is stopped..")
			return
		default:
			if time.Since(updateStamp) > repo.config.UpdateTimeout {
				if err := repo.update(ctx); err != nil {
					repo.log.Err(err).Msg("Failed to update rates")
				}
			}
		}
	}
}

func (repo *PostgresRepo) update(ctx context.Context) error {
	rates, err := fetcher.FetchRates(ctx)
	if err != nil {
		return errors.New("Failed to fetch rates")
	}
	for currencyCode, value := range rates {
		rate := &model.Rate{
			CurrencyCode: currencyCode,
			Value:        value,
		}
		err := repo.setRate(ctx, rate)
		if err != nil {
			return errors.New("Failed to set rate: " + currencyCode)
		}
	}
	return nil
}
