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
		value:       "50051",
		description: "Server port",
	},

	{
		name:        "storage.postgresURL",
		typing:      "string",
		value:       "postgres://postgres:post@localhost:6432/exchange?sslmode=disable",
		description: "Postgres database URL",
	},
	{
		name:        "storage.postgresTimeout",
		typing:      "duration",
		value:       "2s",
		description: "Timeout for database connection",
	},
	{
		name:        "storage.postgresConnectAttempts",
		typing:      "int",
		value:       10,
		description: "Number of database connection attempts",
	},
	{
		name:        "storage.postgresMaxIdleTime",
		typing:      "duration",
		value:       "5m",
		description: "Maximum idle time for database connections",
	},
	{
		name:        "storage.postgresMaxOpenConns",
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
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}
