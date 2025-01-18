package pg

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mizmorr/gw_currency/gw-exchanger/internal/config"
	"github.com/mizmorr/gw_currency/gw-exchanger/pkg/utils/connector"
	logger "github.com/mizmorr/loggerm"
)

type pg struct {
	*pgxpool.Pool
}

var (
	once       sync.Once
	pgInstance *pg
)

func newPg(ctx context.Context) (*pg, error) {
	confg := config.Get()
	log := logger.GetLoggerFromContext(ctx)

	log.Debug().Msg("Database url checking...")
	if confg.PostgresURL == "" {
		return nil, errors.New("No database URL provided")
	}

	log.Debug().Msg("Database config parsing...")
	poolConfig, err := pgxpool.ParseConfig(confg.PostgresURL)
	if err != nil {
		return nil, errors.Wrap(err, "Parse config failed")
	}

	poolConfig.MaxConnIdleTime = confg.PostgresMaxIdleTime
	poolConfig.HealthCheckPeriod = confg.PostgresHealthCheckPeriod
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create pgx pool")
	}

	pgInstance = &pg{pool}

	return pgInstance, nil
}

func dial(ctx context.Context, connectAttempts int, timeout time.Duration) error {
	var err error
	once.Do(func() {
		err = connector.EstablishConnection(ctx, pgInstance.Pool, connectAttempts, timeout)
	})
	return err
}
