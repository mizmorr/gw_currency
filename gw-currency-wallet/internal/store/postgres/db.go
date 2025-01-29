package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mizmorr/gw_currency/gw-currency-wallet/internal/config"
	"github.com/mizmorr/gw_currency/gw-exchanger/pkg/utils/connector"
	logger "github.com/mizmorr/loggerm"
	"github.com/pkg/errors"
)

type db struct {
	*pgxpool.Pool
}

var (
	pgInstance *db
	once       sync.Once
)

func newDBConnector(ctx context.Context, confg *config.Config) (*db, error) {
	log := logger.GetLoggerFromContext(ctx)

	log.Debug().Msg("Database config parsing...")

	poolConfig, err := parseConfig(confg)
	if err != nil {
		log.Err(err).Msg("Error occured while parsing the config")

		return nil, errors.Wrap(err, "Failed to parse database config")
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Err(err).Msg("Error occured while creating pgx pool")

		return nil, errors.Wrap(err, "Failed to create pgx pool")
	}
	pgInstance = &db{pool}

	return pgInstance, nil
}

func parseConfig(confg *config.Config) (*pgxpool.Config, error) {
	if confg.URL == "" {
		return nil, errors.New("No database URL provided")
	}

	poolConfig, err := pgxpool.ParseConfig(confg.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Parse config failed")
	}

	poolConfig.MaxConnIdleTime = confg.MaxIdleTime
	poolConfig.HealthCheckPeriod = confg.HealthCheckPeriod
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	return poolConfig, nil
}

func dial(ctx context.Context, connectAttempts int, timeout time.Duration) error {
	var err error
	once.Do(func() {
		err = connector.EstablishConnection(ctx, pgInstance.Pool, connectAttempts, timeout)
	})
	return err
}

func (d *db) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return d.Pool.QueryRow(ctx, sql, args...)
}

func (d *db) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.Pool.Exec(ctx, sql, args)
}

func (d *db) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return d.Pool.Query(ctx, sql, args...)
}

func (d *db) Close() {
	d.Pool.Close()
}

func (d *db) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return d.Pool.Acquire(ctx)
}

func (d *db) Begin(ctx context.Context) (pgx.Tx, error) {
	return d.Pool.Begin(ctx)
}
