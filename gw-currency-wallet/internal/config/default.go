package config

var defaults = []option{
	{
		name:        "logger.level",
		typing:      "string",
		value:       "debug",
		description: "Level of logging",
	},
	{
		name:        "logger.pathFile",
		typing:      "string",
		value:       "./logs/app.log",
		description: "Path to the log file",
	},
	{
		name:        "listen.host",
		typing:      "string",
		value:       "localhost",
		description: "Server host",
	},
	{
		name:        "listen.port",
		typing:      "string",
		value:       "8080",
		description: "Server port",
	},

	{
		name:        "storage.postgres.URL",
		typing:      "string",
		value:       "postgres://postgres:post@localhost:6432/exchange?sslmode=disable",
		description: "Postgres database URL",
	},
	{
		name:        "storage.postgres.Timeout",
		typing:      "duration",
		value:       "2s",
		description: "Timeout for database connection",
	},
	{
		name:        "storage.postgres.ConnectAttempts",
		typing:      "int",
		value:       10,
		description: "Number of database connection attempts",
	},
	{
		name:        "storage.postgres.MaxIdleTime",
		typing:      "duration",
		value:       "5m",
		description: "Maximum idle time for database connections",
	},
	{
		name:        "storage.postgres.MaxOpenConns",
		typing:      "int",
		value:       100,
		description: "Maximum number of open database connections",
	},
	{
		name:        "storage.postgresHealthCheckPeriod",
		typing:      "duration",
		value:       "5m",
		description: "Period for database health check",
	},
	{
		name:        "workers.keepAliveTimeout",
		typing:      "duration",
		value:       "5s",
		description: "Timeout for worker keep-alive",
	},
	{
		name:        "workers.updateTimeout",
		typing:      "duration",
		value:       "10m",
		description: "Timeout for worker update",
	},
	{
		name:        "storage.redis.host",
		typing:      "string",
		value:       "localhost",
		description: "Redis server host",
	},
	{
		name:        "storage.redis.port",
		typing:      "int",
		value:       6379,
		description: "Redis server port",
	},
	{
		name:        "storage.redis.password",
		typing:      "string",
		value:       "mor",
		description: "Redis server password",
	},
	{
		name:        "storage.redis.db",
		typing:      "int",
		value:       0,
		description: "Redis database index",
	},
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}
