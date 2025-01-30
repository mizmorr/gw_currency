package postgres

import (
	"context"
	"time"

	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	logger "github.com/mizmorr/loggerm"
)

type PostgresRepo struct {
	db     *db
	stop   chan interface{}
	config *config.Config
	log    *logger.Logger
}

func NewPostgresRepo(ctx context.Context) (*PostgresRepo, error) {
	config := config.Get()

	log := logger.GetLoggerFromContext(ctx)

	ch := make(chan interface{})

	db, err := newDBConnector(ctx, config)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{
		stop:   ch,
		config: config,
		log:    log,
		db:     db,
	}, nil
}

func (repo *PostgresRepo) Start(ctx context.Context) error {
	err := dial(ctx, repo.config.ConnectAttempts, repo.config.Timeout)
	if err != nil {
		return err
	}
	go repo.keepAlive(ctx)

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
		err = dial(ctx, repo.config.ConnectAttempts, repo.config.Timeout)
		if err != nil {
			repo.log.Err(err).Msg("Failed to reconnect to PostgreSQL database")
		}
	}
}

func (repo *PostgresRepo) Stop(_ context.Context) error {
	repo.log.Info().Msg("Stopping PostgreSQL repository..")

	repo.stop <- struct{}{}

	repo.db.Close()
	return nil
}
