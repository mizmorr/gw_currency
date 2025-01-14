package config

var defaults = []option{
	{
		name:        "loglevel",
		typing:      "string",
		value:       "info",
		description: "Level of logging",
	},
	{
		name:        "listen.host",
		typing:      "string",
		value:       "localhost",
		description: "Server host",
	},
	{
		name:        "listen.port",
		typing:      "int",
		value:       8080,
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
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}
